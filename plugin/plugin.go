package plugin

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const THRESHOLDS_BREACHED_EXIT_CODE = 99

type shellCommand interface {
	Start() error
	Wait() error
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
	String() string
}

type errorWithExitCode interface {
	ExitCode() int
}

// buildExecCommand returns a shellCommand with the given arguments. The
// return type of shellCommand is for mocking purposes.
func buildExecCommand(name string, args ...string) shellCommand {
	return exec.Command(name, args...)
}

// checkOSStat verifies a file exists at the given path, otherwise returns
// an error
func checkOSStat(path string) error {
	_, err := os.Stat(path)
	return err
}

var (
	validFilePattern = regexp.MustCompile(`^(\./|(\.\./)+)?[a-zA-Z0-9-_/]*[a-zA-Z0-9]\.(json|js)$`)
	// buildCommand can be swapped out for a mock function for unit testing
	buildCommand = buildExecCommand
	// verifyFileExists can be swapped out for a mock function for unit testing
	verifyFileExists = checkOSStat
)

// ConfigFromEnv returns a Config populated with the values of the Vela
// parameters. Script and output paths will be sanitized/validated, and
// an error is returned if the script path is empty or invalid. If the
// output path is invalid, OutputPath is set to ""
func ConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	cfg.ScriptPath = sanitizeFilePath(os.Getenv("PARAMETER_SCRIPT_PATH"))
	cfg.OutputPath = sanitizeFilePath(os.Getenv("PARAMETER_OUTPUT_PATH"))
	cfg.FailOnThresholdBreach = !strings.EqualFold(os.Getenv("PARAMETER_FAIL_ON_THRESHOLD_BREACH"), "false")
	cfg.ProjektorCompatMode = strings.EqualFold(os.Getenv("PARAMETER_PROJEKTOR_COMPAT_MODE"), "true")
	cfg.LogProgress = strings.EqualFold(os.Getenv("PARAMETER_LOG_PROGRESS"), "true")

	if cfg.ScriptPath == "" || !strings.HasSuffix(cfg.ScriptPath, ".js") {
		return nil, fmt.Errorf("invalid script file. provide the filepath to a JavaScript file in plugin parameter 'script_path' (e.g. 'script_path: \"/k6-test/script.js\"'). the filepath must follow the regular expression `^[a-zA-Z0-9-_/]*[a-zA-Z0-9]+\\.(json|js)$`")
	}

	return cfg, nil
}

// sanitizeFilePath returns the input string if it satisfies the pattern
// for a valid filepath, and an empty string otherwise
func sanitizeFilePath(input string) string {
	return validFilePattern.FindString(input)
}

// buildK6Command returns a shellCommand that will execute K6 tests
// using the script path, output path, and output type in cfg
func buildK6Command(cfg *Config) (cmd shellCommand, err error) {
	commandArgs := []string{"run"}
	if !cfg.LogProgress {
		commandArgs = append(commandArgs, "-q")
	}
	if cfg.OutputPath != "" {
		outputDir := filepath.Dir(cfg.OutputPath)
		err = os.MkdirAll(outputDir, os.FileMode(0755))
		if err != nil {
			return
		}

		if cfg.ProjektorCompatMode {
			commandArgs = append(commandArgs, fmt.Sprintf("--summary-export=%s", cfg.OutputPath))
		} else {
			commandArgs = append(commandArgs, "--out", fmt.Sprintf("json=%s", cfg.OutputPath))
		}
	}
	commandArgs = append(commandArgs, cfg.ScriptPath)
	cmd = buildCommand("k6", commandArgs...)
	return
}

// RunPerfTests runs the K6 performance test script located at the
// cfg.ScriptPath and saves the output to cfg.OutputPath if it is present
// and a valid filepath
func RunPerfTests(cfg *Config) error {

	err := verifyFileExists(cfg.ScriptPath)
	if err != nil {
		return fmt.Errorf("Error reading script file at %s: %s", cfg.ScriptPath, err)
	}

	cmd, err := buildK6Command(cfg)
	if err != nil {
		return fmt.Errorf("Error creating output directory: %s", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error getting stdout pipe: %s", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error getting stderr pipe: %s", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %s", err)
	}

	log.Println("Running tests...")

	wg := sync.WaitGroup{}
	wg.Add(2)
	go readLinesFromPipe(stdout, &wg)
	go readLinesFromPipe(stderr, &wg)
	wg.Wait()

	execError := cmd.Wait()

	if execError != nil {
		exitError, ok := execError.(errorWithExitCode)
		if ok && exitError.ExitCode() == THRESHOLDS_BREACHED_EXIT_CODE {
			if cfg.FailOnThresholdBreach {
				return fmt.Errorf("thresholds breached")
			}
		} else {
			return execError
		}
	}

	if cfg.OutputPath != "" {
		path, err := filepath.Abs(cfg.OutputPath)
		if err != nil {
			log.Printf("Error saving output to %s: %s\n", cfg.OutputPath, err)
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

type Config struct {
	ScriptPath            string
	OutputPath            string
	FailOnThresholdBreach bool
	ProjektorCompatMode   bool
	LogProgress           bool
}
