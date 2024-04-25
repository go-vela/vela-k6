package plugin

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/go-vela/vela-k6/types"
)

const thresholdsBreachedExitCode = 99

type pluginType struct {
	config           config
	buildCommand     func(name string, args ...string) types.ShellCommand // buildCommand can be swapped out for a mock function for unit testing.
	verifyFileExists func(path string) error                              // verifyFileExists can be swapped out for a mock function for unit testing.
}

type Plugin interface {
	ConfigFromEnv() error
	RunSetupScript() error
	RunPerfTests() error
}

func New() Plugin {
	return &pluginType{
		buildCommand:     buildExecCommand,
		verifyFileExists: checkOSStat,
	}
}

// buildExecCommand returns a ShellCommand with the given arguments. The
// return type of ShellCommand is for mocking purposes.
func buildExecCommand(name string, args ...string) types.ShellCommand {
	return exec.Command(name, args...)
}

// checkOSStat verifies a file exists at the given path, otherwise returns
// an error.
func checkOSStat(path string) error {
	_, err := os.Stat(path)
	return err
}

var (
	validJSFilePattern    = regexp.MustCompile(`^(\./|(\.\./)+)?[a-zA-Z0-9-_/]*[a-zA-Z0-9]\.js$`)
	validJSONFilePattern  = regexp.MustCompile(`^(\./|(\.\./)+)?[a-zA-Z0-9-_/]*[a-zA-Z0-9]\.json$`)
	validShellFilePattern = regexp.MustCompile(`^(\./|(\.\./)+)?[a-zA-Z0-9-_/]*[a-zA-Z0-9]\.sh$`)
)

// ConfigFromEnv returns a Config populated with the values of the Vela
// parameters. Script and output paths will be sanitized/validated, and
// an error is returned if the script path is empty or invalid. If the
// output path is invalid, OutputPath is set to "".
func (p *pluginType) ConfigFromEnv() error {
	p.config.ScriptPath = sanitizeScriptPath(os.Getenv("PARAMETER_SCRIPT_PATH"))
	p.config.OutputPath = sanitizeOutputPath(os.Getenv("PARAMETER_OUTPUT_PATH"))
	p.config.SetupScriptPath = sanitizeSetupPath(os.Getenv("PARAMETER_SETUP_SCRIPT_PATH"))
	p.config.FailOnThresholdBreach = !strings.EqualFold(os.Getenv("PARAMETER_FAIL_ON_THRESHOLD_BREACH"), "false")
	p.config.ProjektorCompatMode = strings.EqualFold(os.Getenv("PARAMETER_PROJEKTOR_COMPAT_MODE"), "true")
	p.config.LogProgress = strings.EqualFold(os.Getenv("PARAMETER_LOG_PROGRESS"), "true")

	if p.config.ScriptPath == "" || !strings.HasSuffix(p.config.ScriptPath, ".js") {
		p.config = config{} // reset config
		return fmt.Errorf("invalid script file. provide the filepath to a JavaScript file in plugin parameter 'script_path' (e.g. 'script_path: \"/k6-test/script.js\"'). the filepath must follow the regular expression `%s`", validJSFilePattern)
	}

	return nil
}

// sanitizeScriptPath returns the input string if it satisfies the pattern
// for a valid JS filepath, and an empty string otherwise.
func sanitizeScriptPath(input string) string {
	return validJSFilePattern.FindString(input)
}

// sanitizeOutputPath returns the input string if it satisfies the pattern
// for a valid JSON filepath, and an empty string otherwise.
func sanitizeOutputPath(input string) string {
	return validJSONFilePattern.FindString(input)
}

// sanitizeSetupPath returns the input string if it satisfies the pattern
// for a valid .sh filepath, and an empty string otherwise.
func sanitizeSetupPath(input string) string {
	return validShellFilePattern.FindString(input)
}

// buildK6Command returns a ShellCommand that will execute K6 tests
// using the script path, output path, and output type in cfg.
func (p *pluginType) buildK6Command() (cmd types.ShellCommand, err error) {
	commandArgs := []string{"run"}
	if !p.config.LogProgress {
		commandArgs = append(commandArgs, "-q")
	}

	if p.config.OutputPath != "" {
		outputDir := filepath.Dir(p.config.OutputPath)
		if err = os.MkdirAll(outputDir, os.FileMode(0755)); err != nil {
			return
		}

		if p.config.ProjektorCompatMode {
			commandArgs = append(commandArgs, fmt.Sprintf("--summary-export=%s", p.config.OutputPath))
		} else {
			commandArgs = append(commandArgs, "--out", fmt.Sprintf("json=%s", p.config.OutputPath))
		}
	}

	commandArgs = append(commandArgs, p.config.ScriptPath)
	cmd = p.buildCommand("k6", commandArgs...)

	return
}

// RunSetupScript runs the setup script located at the cfg.SetupScriptPath
// if the path is not empty.
func (p *pluginType) RunSetupScript() error {
	if p.config.SetupScriptPath == "" {
		log.Println("No setup script specified, skipping.")
		return nil
	}

	err := p.verifyFileExists(p.config.SetupScriptPath)
	if err != nil {
		return fmt.Errorf("read setup script file at %s: %w", p.config.SetupScriptPath, err)
	}

	cmd := p.buildCommand(p.config.SetupScriptPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	log.Println("Running setup script...")

	wg := sync.WaitGroup{}
	wg.Add(2)

	go readLinesFromPipe(stdout, &wg)
	go readLinesFromPipe(stderr, &wg)
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("run setup script: %w", err)
	}

	return nil
}

// RunPerfTests runs the K6 performance test script located at the
// p.config.ScriptPath and saves the output to p.config.OutputPath if it is present
// and a valid filepath.
func (p *pluginType) RunPerfTests() error {
	err := p.verifyFileExists(p.config.ScriptPath)
	if err != nil {
		return fmt.Errorf("read script file at %s: %w", p.config.ScriptPath, err)
	}

	cmd, err := p.buildK6Command()
	if err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	log.Println("Running tests...")

	wg := sync.WaitGroup{}
	wg.Add(2)

	go readLinesFromPipe(stdout, &wg)
	go readLinesFromPipe(stderr, &wg)
	wg.Wait()

	execError := cmd.Wait()

	if execError != nil {
		var exitError types.ErrorWithExitCode
		ok := errors.As(execError, &exitError)

		if ok && exitError.ExitCode() == thresholdsBreachedExitCode {
			if p.config.FailOnThresholdBreach {
				return fmt.Errorf("thresholds breached")
			}
		} else {
			return execError
		}
	}

	if p.config.OutputPath != "" {
		path, err := filepath.Abs(p.config.OutputPath)
		if err != nil {
			log.Printf("save output to %s: %s\n", p.config.OutputPath, err)
		} else {
			log.Printf("Output file saved at %s\n", path)
		}
	}

	return nil
}

// readLinesFromPipe will read each line from pipe and log it. A WaitGroup
// may optionally be passed in, in which case Done() will be called
// once the pipe is closed.
func readLinesFromPipe(pipe io.ReadCloser, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}
}

type config struct {
	ScriptPath            string
	OutputPath            string
	SetupScriptPath       string
	FailOnThresholdBreach bool
	ProjektorCompatMode   bool
	LogProgress           bool
}
