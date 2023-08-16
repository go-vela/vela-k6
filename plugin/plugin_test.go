package plugin

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"testing"

	"github.com/go-vela/vela-k6/plugin/mock"
	"github.com/stretchr/testify/assert"
)

func setFilePathEnvs(t *testing.T) {
	t.Setenv("PARAMETER_SCRIPT_PATH", "./test/script.js")
	t.Setenv("PARAMETER_OUTPUT_PATH", "./output.json")
	t.Setenv("PARAMETER_SETUP_SCRIPT_PATH", "./test/setup.sh")
}

func clearEnvironment(t *testing.T) {
	t.Setenv("PARAMETER_SCRIPT_PATH", "")
	t.Setenv("PARAMETER_OUTPUT_PATH", "")
	t.Setenv("PARAMETER_SETUP_SCRIPT_PATH", "")
	t.Setenv("PARAMETER_PROJEKTOR_COMPAT_MODE", "")
	t.Setenv("PARAMETER_FAIL_ON_THRESHOLD_BREACH", "")
	t.Setenv("PARAMETER_LOG_PROGRESS", "")
}

func TestSanitizeScriptPath(t *testing.T) {
	t.Run("Valid Filepaths", func(t *testing.T) {
		assert.Equal(t, "file.js", sanitizeScriptPath("file.js"))
		assert.Equal(t, "./file.js", sanitizeScriptPath("./file.js"))
		assert.Equal(t, "../file.js", sanitizeScriptPath("../file.js"))
		assert.Equal(t, "../../../file.js", sanitizeScriptPath("../../../file.js"))
		assert.Equal(t, "file-dash_underscore.js", sanitizeScriptPath("file-dash_underscore.js"))
		assert.Equal(t, "path/to/file.js", sanitizeScriptPath("path/to/file.js"))
		assert.Equal(t, "/path/to/file.js", sanitizeScriptPath("/path/to/file.js"))
	})

	t.Run("Invalid Filepaths", func(t *testing.T) {
		assert.Equal(t, "", sanitizeScriptPath(".../file.js"))
		assert.Equal(t, "", sanitizeScriptPath("./../file.js"))
		assert.Equal(t, "", sanitizeScriptPath("*/file.js"))
		assert.Equal(t, "", sanitizeScriptPath(".file.js"))
		assert.Equal(t, "", sanitizeScriptPath("/.js"))
		assert.Equal(t, "", sanitizeScriptPath("-.js"))
		assert.Equal(t, "", sanitizeScriptPath("_.js"))
		assert.Equal(t, "", sanitizeScriptPath("_invalid$name.js"))
		assert.Equal(t, "", sanitizeScriptPath("invalid$name.js"))
		assert.Equal(t, "", sanitizeScriptPath("invalidformat.png"))
		assert.Equal(t, "", sanitizeScriptPath("file.js; rm -rf /"))
		assert.Equal(t, "", sanitizeScriptPath("file.js && suspicious-call"))
	})
}

func TestSanitizeOutputPath(t *testing.T) {
	t.Run("Valid Filepaths", func(t *testing.T) {
		assert.Equal(t, "file.json", sanitizeOutputPath("file.json"))
		assert.Equal(t, "./file.json", sanitizeOutputPath("./file.json"))
		assert.Equal(t, "../file.json", sanitizeOutputPath("../file.json"))
		assert.Equal(t, "../../../file.json", sanitizeOutputPath("../../../file.json"))
		assert.Equal(t, "file-dash_underscore.json", sanitizeOutputPath("file-dash_underscore.json"))
		assert.Equal(t, "path/to/file.json", sanitizeOutputPath("path/to/file.json"))
		assert.Equal(t, "/path/to/file.json", sanitizeOutputPath("/path/to/file.json"))
	})

	t.Run("Invalid Filepaths", func(t *testing.T) {
		assert.Equal(t, "", sanitizeOutputPath(".../file.json"))
		assert.Equal(t, "", sanitizeOutputPath("./../file.json"))
		assert.Equal(t, "", sanitizeOutputPath("*/file.json"))
		assert.Equal(t, "", sanitizeOutputPath(".file.json"))
		assert.Equal(t, "", sanitizeOutputPath("/.json"))
		assert.Equal(t, "", sanitizeOutputPath("-.json"))
		assert.Equal(t, "", sanitizeOutputPath("_.json"))
		assert.Equal(t, "", sanitizeOutputPath("_invalid$name.json"))
		assert.Equal(t, "", sanitizeOutputPath("invalid$name.json"))
		assert.Equal(t, "", sanitizeOutputPath("invalidformat.png"))
		assert.Equal(t, "", sanitizeOutputPath("file.json; rm -rf /"))
		assert.Equal(t, "", sanitizeOutputPath("file.json && suspicious-call"))
	})
}

func TestSanitizeSetupPath(t *testing.T) {
	t.Run("Valid Filepaths", func(t *testing.T) {
		assert.Equal(t, "file.sh", sanitizeSetupPath("file.sh"))
		assert.Equal(t, "./file.sh", sanitizeSetupPath("./file.sh"))
		assert.Equal(t, "../file.sh", sanitizeSetupPath("../file.sh"))
		assert.Equal(t, "../../../file.sh", sanitizeSetupPath("../../../file.sh"))
		assert.Equal(t, "file-dash_underscore.sh", sanitizeSetupPath("file-dash_underscore.sh"))
		assert.Equal(t, "path/to/file.sh", sanitizeSetupPath("path/to/file.sh"))
		assert.Equal(t, "/path/to/file.sh", sanitizeSetupPath("/path/to/file.sh"))
	})

	t.Run("Invalid Filepaths", func(t *testing.T) {
		assert.Equal(t, "", sanitizeSetupPath(".../file.sh"))
		assert.Equal(t, "", sanitizeSetupPath("./../file.sh"))
		assert.Equal(t, "", sanitizeSetupPath("*/file.sh"))
		assert.Equal(t, "", sanitizeSetupPath(".file.sh"))
		assert.Equal(t, "", sanitizeSetupPath("/.sh"))
		assert.Equal(t, "", sanitizeSetupPath("-.sh"))
		assert.Equal(t, "", sanitizeSetupPath("_.sh"))
		assert.Equal(t, "", sanitizeSetupPath("_invalid$name.sh"))
		assert.Equal(t, "", sanitizeSetupPath("invalid$name.sh"))
		assert.Equal(t, "", sanitizeSetupPath("invalidformat.png"))
		assert.Equal(t, "", sanitizeSetupPath("file.sh; rm -rf /"))
		assert.Equal(t, "", sanitizeSetupPath("file.sh && suspicious-call"))
	})
}

func TestConfigFromEnv(t *testing.T) {
	clearEnvironment(t)
	t.Run("Files Only", func(t *testing.T) {
		setFilePathEnvs(t)
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		assert.Equal(t, "./test/script.js", cfg.ScriptPath)
		assert.Equal(t, "./output.json", cfg.OutputPath)
	})
	t.Run("Non-Default Options", func(t *testing.T) {
		setFilePathEnvs(t)
		t.Setenv("PARAMETER_PROJEKTOR_COMPAT_MODE", "true")
		t.Setenv("PARAMETER_FAIL_ON_THRESHOLD_BREACH", "false")

		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		assert.Equal(t, "./test/script.js", cfg.ScriptPath)
		assert.Equal(t, "./output.json", cfg.OutputPath)
		assert.True(t, cfg.ProjektorCompatMode)
		assert.False(t, cfg.FailOnThresholdBreach)
	})
	t.Run("Invalid Script Path", func(t *testing.T) {
		t.Setenv("PARAMETER_SCRIPT_PATH", "./script.png")
		cfg, err := ConfigFromEnv()
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}

func TestBuildK6Command(t *testing.T) {
	clearEnvironment(t)
	t.Run("No Output", func(t *testing.T) {
		t.Setenv("PARAMETER_SCRIPT_PATH", "./test/script.js")
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		cmd, err := buildK6Command(cfg)
		assert.NoError(t, err)
		assert.Contains(t, cmd.String(), "k6 run -q ./test/script.js")
	})
	t.Run("Projektor Compat Output", func(t *testing.T) {
		setFilePathEnvs(t)
		t.Setenv("PARAMETER_PROJEKTOR_COMPAT_MODE", "true")
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		cmd, err := buildK6Command(cfg)
		assert.NoError(t, err)
		assert.Contains(t, cmd.String(), "k6 run -q --summary-export=./output.json ./test/script.js")
	})
	t.Run("K6 Recommended Output", func(t *testing.T) {
		setFilePathEnvs(t)
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		cmd, err := buildK6Command(cfg)
		assert.NoError(t, err)
		assert.Contains(t, cmd.String(), "k6 run -q --out json=./output.json ./test/script.js")
	})
	t.Run("Verbose logging", func(t *testing.T) {
		t.Setenv("PARAMETER_SCRIPT_PATH", "./test/script.js")
		t.Setenv("PARAMETER_LOG_PROGRESS", "true")
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		cmd, err := buildK6Command(cfg)
		assert.NoError(t, err)
		assert.Contains(t, cmd.String(), "k6 run ./test/script.js")
	})
}

func TestRunSetupScript(t *testing.T) {
	clearEnvironment(t)

	buildCommand = mock.CommandBuilderWithError(nil)
	verifyFileExists = func(path string) error {
		if path != "./test/setup.sh" {
			return fmt.Errorf("File does not exist at path %s", path)
		}

		return nil
	}

	defer func() {
		buildCommand = buildExecCommand
		verifyFileExists = checkOSStat
	}()

	t.Run("Successful setup script", func(t *testing.T) {
		setFilePathEnvs(t)
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunSetupScript(cfg)
		assert.NoError(t, err)
	})

	t.Run("No setup script", func(t *testing.T) {
		setFilePathEnvs(t)
		t.Setenv("PARAMETER_SETUP_SCRIPT_PATH", "")
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunSetupScript(cfg)
		assert.NoError(t, err)
	})

	t.Run("Script file not present", func(t *testing.T) {
		setFilePathEnvs(t)
		t.Setenv("PARAMETER_SETUP_SCRIPT_PATH", "./test/doesnotexist.sh")
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunSetupScript(cfg)
		assert.ErrorContains(t, err, "read setup script file at")
	})

	t.Run("Setup script exec error", func(t *testing.T) {
		buildCommand = mock.CommandBuilderWithError(fmt.Errorf("some setup error"))
		setFilePathEnvs(t)
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunSetupScript(cfg)
		assert.ErrorContains(t, err, "run setup script: some setup error")
	})
}

func TestRunPerfTests(t *testing.T) {
	clearEnvironment(t)

	buildCommand = mock.CommandBuilderWithError(nil)
	verifyFileExists = func(path string) error {
		if path != "./test/script.js" {
			return fmt.Errorf("File does not exist at path %s", path)
		}

		return nil
	}

	defer func() {
		buildCommand = buildExecCommand
		verifyFileExists = checkOSStat
	}()
	t.Run("Successful Perf Tests", func(t *testing.T) {
		setFilePathEnvs(t)
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunPerfTests(cfg)
		assert.NoError(t, err)
	})

	t.Run("Script file not present", func(t *testing.T) {
		setFilePathEnvs(t)
		t.Setenv("PARAMETER_SCRIPT_PATH", "./test/doesnotexist.js")
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunPerfTests(cfg)
		assert.ErrorContains(t, err, "read script file at")
	})

	t.Run("Error if thresholds breached", func(t *testing.T) {
		buildCommand = mock.CommandBuilderWithError(&mock.ThresholdError{})
		setFilePathEnvs(t)
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunPerfTests(cfg)
		assert.ErrorContains(t, err, "thresholds breached")
	})

	t.Run("No error if thresholds breached", func(t *testing.T) {
		buildCommand = mock.CommandBuilderWithError(&mock.ThresholdError{})
		setFilePathEnvs(t)
		t.Setenv("PARAMETER_FAIL_ON_THRESHOLD_BREACH", "false")
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunPerfTests(cfg)
		assert.NoError(t, err)
	})

	t.Run("Other exec error", func(t *testing.T) {
		buildCommand = mock.CommandBuilderWithError(fmt.Errorf("some exec error"))
		setFilePathEnvs(t)
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunPerfTests(cfg)
		assert.ErrorContains(t, err, "some exec error")
	})
}

func TestReadLinesFromPipe(t *testing.T) {
	t.Run("Reads from pipe and closes", func(t *testing.T) {
		var buf bytes.Buffer
		prevOut := log.Writer()
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(prevOut)
		}()
		line1 := "this is line 1"
		line2 := "this is line 2"
		reader := io.NopCloser(strings.NewReader(fmt.Sprintf("%s\n%s", line1, line2)))

		// same wait group logic as used in plugin
		wg := sync.WaitGroup{}
		wg.Add(1)
		go readLinesFromPipe(reader, &wg)
		wg.Wait()

		logLine, err := buf.ReadString('\n')
		assert.NoError(t, err)
		assert.Contains(t, logLine, line1)
		logLine, err = buf.ReadString('\n')
		assert.NoError(t, err)
		assert.Contains(t, logLine, line2)
	})
}
