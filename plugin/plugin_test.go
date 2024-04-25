package plugin

import (
	"bytes"
	"errors"
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
		t.Parallel()

		assert.Equal(t, "file.js", sanitizeScriptPath("file.js"))
		assert.Equal(t, "./file.js", sanitizeScriptPath("./file.js"))
		assert.Equal(t, "../file.js", sanitizeScriptPath("../file.js"))
		assert.Equal(t, "../../../file.js", sanitizeScriptPath("../../../file.js"))
		assert.Equal(t, "file-dash_underscore.js", sanitizeScriptPath("file-dash_underscore.js"))
		assert.Equal(t, "path/to/file.js", sanitizeScriptPath("path/to/file.js"))
		assert.Equal(t, "/path/to/file.js", sanitizeScriptPath("/path/to/file.js"))
	})

	t.Run("Invalid Filepaths", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

		assert.Equal(t, "file.json", sanitizeOutputPath("file.json"))
		assert.Equal(t, "./file.json", sanitizeOutputPath("./file.json"))
		assert.Equal(t, "../file.json", sanitizeOutputPath("../file.json"))
		assert.Equal(t, "../../../file.json", sanitizeOutputPath("../../../file.json"))
		assert.Equal(t, "file-dash_underscore.json", sanitizeOutputPath("file-dash_underscore.json"))
		assert.Equal(t, "path/to/file.json", sanitizeOutputPath("path/to/file.json"))
		assert.Equal(t, "/path/to/file.json", sanitizeOutputPath("/path/to/file.json"))
	})

	t.Run("Invalid Filepaths", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

		assert.Equal(t, "file.sh", sanitizeSetupPath("file.sh"))
		assert.Equal(t, "./file.sh", sanitizeSetupPath("./file.sh"))
		assert.Equal(t, "../file.sh", sanitizeSetupPath("../file.sh"))
		assert.Equal(t, "../../../file.sh", sanitizeSetupPath("../../../file.sh"))
		assert.Equal(t, "file-dash_underscore.sh", sanitizeSetupPath("file-dash_underscore.sh"))
		assert.Equal(t, "path/to/file.sh", sanitizeSetupPath("path/to/file.sh"))
		assert.Equal(t, "/path/to/file.sh", sanitizeSetupPath("/path/to/file.sh"))
	})

	t.Run("Invalid Filepaths", func(t *testing.T) {
		t.Parallel()

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

		p := &pluginType{}
		err := p.ConfigFromEnv()
		assert.NoError(t, err)
		assert.Equal(t, "./test/script.js", p.config.ScriptPath)
		assert.Equal(t, "./output.json", p.config.OutputPath)
	})
	t.Run("Non-Default Options", func(t *testing.T) {
		setFilePathEnvs(t)
		t.Setenv("PARAMETER_PROJEKTOR_COMPAT_MODE", "true")
		t.Setenv("PARAMETER_FAIL_ON_THRESHOLD_BREACH", "false")

		p := &pluginType{}
		err := p.ConfigFromEnv()
		assert.NoError(t, err)
		assert.Equal(t, "./test/script.js", p.config.ScriptPath)
		assert.Equal(t, "./output.json", p.config.OutputPath)
		assert.True(t, p.config.ProjektorCompatMode)
		assert.False(t, p.config.FailOnThresholdBreach)
	})
	t.Run("Invalid Script Path", func(t *testing.T) {
		t.Setenv("PARAMETER_SCRIPT_PATH", "./script.png")

		p := &pluginType{}
		err := p.ConfigFromEnv()
		assert.Error(t, err)
		assert.Empty(t, p.config)
	})
}

func TestBuildK6Command(t *testing.T) {
	t.Run("No Output", func(t *testing.T) {
		t.Parallel()
		p := &pluginType{
			config:           config{ScriptPath: "./test/script.js"},
			buildCommand:     buildExecCommand,
			verifyFileExists: checkOSStat,
		}

		cmd, err := p.buildK6Command()
		assert.NoError(t, err)
		assert.Contains(t, cmd.String(), "k6 run -q ./test/script.js")
	})
	t.Run("Projektor Compat Output", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:          "./test/script.js",
				OutputPath:          "./output.json",
				SetupScriptPath:     "./test/setup.sh",
				ProjektorCompatMode: true,
			},
			buildCommand:     buildExecCommand,
			verifyFileExists: checkOSStat,
		}

		cmd, err := p.buildK6Command()
		assert.NoError(t, err)
		assert.Contains(t, cmd.String(), "k6 run -q --summary-export=./output.json ./test/script.js")
	})
	t.Run("K6 Recommended Output", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:      "./test/script.js",
				OutputPath:      "./output.json",
				SetupScriptPath: "./test/setup.sh",
			},
			buildCommand:     buildExecCommand,
			verifyFileExists: checkOSStat,
		}

		cmd, err := p.buildK6Command()
		assert.NoError(t, err)
		assert.Contains(t, cmd.String(), "k6 run -q --out json=./output.json ./test/script.js")
	})
	t.Run("Verbose logging", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:  "./test/script.js",
				LogProgress: true,
			},
			buildCommand:     buildExecCommand,
			verifyFileExists: checkOSStat,
		}

		cmd, err := p.buildK6Command()
		assert.NoError(t, err)
		assert.Contains(t, cmd.String(), "k6 run ./test/script.js")
	})
}

func TestRunSetupScript(t *testing.T) {
	buildCommand := mock.CommandBuilderWithError(nil, nil, nil, nil)
	verifyFileExists := func(path string) error {
		if path != "./test/setup.sh" {
			return fmt.Errorf("File does not exist at path %s", path)
		}

		return nil
	}

	t.Run("Successful setup script", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:      "./test/script.js",
				OutputPath:      "./output.json",
				SetupScriptPath: "./test/setup.sh",
			},
			buildCommand:     buildCommand,
			verifyFileExists: verifyFileExists,
		}

		err := p.RunSetupScript()
		assert.NoError(t, err)
	})

	t.Run("No setup script", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath: "./test/script.js",
				OutputPath: "./output.json",
			},
			buildCommand:     buildCommand,
			verifyFileExists: verifyFileExists,
		}
		err := p.RunSetupScript()
		assert.NoError(t, err)
	})

	t.Run("Script file not present", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:      "./test/script.js",
				OutputPath:      "./output.json",
				SetupScriptPath: "./test/doesnotexist.sh",
			},
			buildCommand:     buildCommand,
			verifyFileExists: verifyFileExists,
		}

		err := p.RunSetupScript()
		assert.ErrorContains(t, err, "read setup script file at")
	})
	t.Run("StdoutPipe error", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:      "./test/script.js",
				OutputPath:      "./output.json",
				SetupScriptPath: "./test/setup.sh",
			},
			buildCommand:     mock.CommandBuilderWithError(nil, errors.New("some error"), nil, nil),
			verifyFileExists: verifyFileExists,
		}

		err := p.RunSetupScript()
		assert.ErrorContains(t, err, "get stdout pipe")
	})
	t.Run("StderrPipeErr error", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:      "./test/script.js",
				OutputPath:      "./output.json",
				SetupScriptPath: "./test/setup.sh",
			},
			buildCommand:     mock.CommandBuilderWithError(nil, nil, errors.New("some error"), nil),
			verifyFileExists: verifyFileExists,
		}

		err := p.RunSetupScript()
		assert.ErrorContains(t, err, "get stderr pipe")
	})
	t.Run("Start error", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:      "./test/script.js",
				OutputPath:      "./output.json",
				SetupScriptPath: "./test/setup.sh",
			},
			buildCommand:     mock.CommandBuilderWithError(nil, nil, nil, errors.New("some error")),
			verifyFileExists: verifyFileExists,
		}

		err := p.RunSetupScript()
		assert.ErrorContains(t, err, "start command")
	})
	t.Run("Setup script exec error", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:      "./test/script.js",
				OutputPath:      "./output.json",
				SetupScriptPath: "./test/setup.sh",
			},
			buildCommand:     mock.CommandBuilderWithError(errors.New("some setup error"), nil, nil, nil),
			verifyFileExists: verifyFileExists,
		}

		err := p.RunSetupScript()
		assert.ErrorContains(t, err, "run setup script: some setup error")
	})
}

func TestRunPerfTests(t *testing.T) {
	buildCommand := mock.CommandBuilderWithError(nil, nil, nil, nil)
	verifyFileExists := func(path string) error {
		if path != "./test/script.js" {
			return fmt.Errorf("File does not exist at path %s", path)
		}

		return nil
	}

	t.Run("Successful Perf Tests", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:      "./test/script.js",
				OutputPath:      "./output.json",
				SetupScriptPath: "./test/setup.sh",
			},
			buildCommand:     buildCommand,
			verifyFileExists: verifyFileExists,
		}
		assert.NoError(t, p.RunPerfTests())
	})

	t.Run("Script file not present", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:      "./test/doesnotexist.js",
				OutputPath:      "./output.json",
				SetupScriptPath: "./test/setup.sh",
			},
			buildCommand:     buildCommand,
			verifyFileExists: verifyFileExists,
		}
		assert.ErrorContains(t, p.RunPerfTests(), "read script file at")
	})

	t.Run("Error if thresholds breached", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:            "./test/script.js",
				OutputPath:            "./output.json",
				SetupScriptPath:       "./test/setup.sh",
				FailOnThresholdBreach: true,
			},
			buildCommand:     mock.CommandBuilderWithError(&mock.ThresholdError{}, nil, nil, nil),
			verifyFileExists: verifyFileExists,
		}
		assert.ErrorContains(t, p.RunPerfTests(), "thresholds breached")
	})

	t.Run("No error if thresholds breached", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:            "./test/script.js",
				OutputPath:            "./output.json",
				SetupScriptPath:       "./test/setup.sh",
				FailOnThresholdBreach: false,
			},
			buildCommand:     mock.CommandBuilderWithError(&mock.ThresholdError{}, nil, nil, nil),
			verifyFileExists: verifyFileExists,
		}

		assert.NoError(t, p.RunPerfTests())
	})

	t.Run("Other exec error", func(t *testing.T) {
		t.Parallel()

		p := &pluginType{
			config: config{
				ScriptPath:            "./test/script.js",
				OutputPath:            "./output.json",
				SetupScriptPath:       "./test/setup.sh",
				FailOnThresholdBreach: true,
			},
			buildCommand:     mock.CommandBuilderWithError(errors.New("some exec error"), nil, nil, nil),
			verifyFileExists: verifyFileExists,
		}
		assert.ErrorContains(t, p.RunPerfTests(), "some exec error")
	})
}

func TestReadLinesFromPipe(t *testing.T) {
	t.Run("Reads from pipe and closes", func(t *testing.T) {
		var buf bytes.Buffer

		log.SetOutput(&buf)

		prevOut := log.Writer()
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

func TestCheckOSStat(t *testing.T) {
	assert.Error(t, checkOSStat("./test/doesnotexist.js"))
	assert.NoError(t, checkOSStat("plugin.go"))
}

func TestNew(t *testing.T) {
	p := New()
	assert.NotNil(t, p)
}
