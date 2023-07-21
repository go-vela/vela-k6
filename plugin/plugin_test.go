// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package plugin

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setFilePathEnvs() {
	os.Setenv("PARAMETER_SCRIPT_PATH", "./test/script.js")
	os.Setenv("PARAMETER_OUTPUT_PATH", "./output.json")
}

func TestSanitizeFilePath(t *testing.T) {
	t.Run("Valid Filepaths", func(t *testing.T) {
		assert.Equal(t, "file.js", sanitizeFilePath("file.js"))
		assert.Equal(t, "./file.js", sanitizeFilePath("./file.js"))
		assert.Equal(t, "../file.js", sanitizeFilePath("../file.js"))
		assert.Equal(t, "../../../file.js", sanitizeFilePath("../../../file.js"))
		assert.Equal(t, "file.json", sanitizeFilePath("file.json"))
		assert.Equal(t, "file-dash_underscore.json", sanitizeFilePath("file-dash_underscore.json"))
		assert.Equal(t, "path/to/file.js", sanitizeFilePath("path/to/file.js"))
		assert.Equal(t, "/path/to/file.js", sanitizeFilePath("/path/to/file.js"))
		assert.Equal(t, "path/to/file.json", sanitizeFilePath("path/to/file.json"))
	})

	t.Run("Invalid Filepaths", func(t *testing.T) {
		assert.Equal(t, "", sanitizeFilePath(".../file.js"))
		assert.Equal(t, "", sanitizeFilePath("./../file.js"))
		assert.Equal(t, "", sanitizeFilePath("*/file.js"))
		assert.Equal(t, "", sanitizeFilePath(".file.js"))
		assert.Equal(t, "", sanitizeFilePath("/.json"))
		assert.Equal(t, "", sanitizeFilePath("-.json"))
		assert.Equal(t, "", sanitizeFilePath("_.json"))
		assert.Equal(t, "", sanitizeFilePath("_invalid$name.json"))
		assert.Equal(t, "", sanitizeFilePath("invalid$name.js"))
		assert.Equal(t, "", sanitizeFilePath("invalidformat.png"))
		assert.Equal(t, "", sanitizeFilePath("file.js; rm -rf /"))
		assert.Equal(t, "", sanitizeFilePath("file.js && suspicious-call"))
	})
}

func TestConfigFromEnv(t *testing.T) {
	t.Run("Files Only", func(t *testing.T) {
		setFilePathEnvs()
		defer os.Clearenv()
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		assert.Equal(t, "./test/script.js", cfg.ScriptPath)
		assert.Equal(t, "./output.json", cfg.OutputPath)
	})
	t.Run("Non-Default Options", func(t *testing.T) {
		setFilePathEnvs()
		os.Setenv("PARAMETER_PROJEKTOR_COMPAT_MODE", "true")
		os.Setenv("PARAMETER_FAIL_ON_THRESHOLD_BREACH", "false")
		defer os.Clearenv()
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		assert.Equal(t, "./test/script.js", cfg.ScriptPath)
		assert.Equal(t, "./output.json", cfg.OutputPath)
		assert.True(t, cfg.ProjektorCompatMode)
		assert.False(t, cfg.FailOnThresholdBreach)
	})
	t.Run("Invalid Script Path", func(t *testing.T) {
		os.Setenv("PARAMETER_SCRIPT_PATH", "./script.png")
		defer os.Clearenv()
		cfg, err := ConfigFromEnv()
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}

func TestBuildK6Command(t *testing.T) {
	t.Run("No Output", func(t *testing.T) {
		os.Setenv("PARAMETER_SCRIPT_PATH", "./test/script.js")
		cfg, err := ConfigFromEnv()
		defer os.Clearenv()
		assert.NoError(t, err)
		cmd, err := buildK6Command(cfg)
		assert.NoError(t, err)
		assert.Contains(t, cmd.String(), "k6 run -q ./test/script.js")
	})
	t.Run("Projektor Compat Output", func(t *testing.T) {
		setFilePathEnvs()
		os.Setenv("PARAMETER_PROJEKTOR_COMPAT_MODE", "true")
		defer os.Clearenv()
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		cmd, err := buildK6Command(cfg)
		assert.NoError(t, err)
		assert.Contains(t, cmd.String(), "k6 run -q --summary-export=./output.json ./test/script.js")
	})
	t.Run("K6 Recommended Output", func(t *testing.T) {
		setFilePathEnvs()
		defer os.Clearenv()
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		cmd, err := buildK6Command(cfg)
		assert.NoError(t, err)
		assert.Contains(t, cmd.String(), "k6 run -q --out json=./output.json ./test/script.js")
	})
}

func TestRunPerfTests(t *testing.T) {
	buildCommand = MockCommandBuilderWithError(nil)
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
		setFilePathEnvs()
		defer os.Clearenv()
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunPerfTests(cfg)
		assert.NoError(t, err)
	})

	t.Run("Error if thresholds breached", func(t *testing.T) {
		buildCommand = MockCommandBuilderWithError(&MockThresholdError{})
		setFilePathEnvs()
		defer os.Clearenv()
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunPerfTests(cfg)
		assert.ErrorContains(t, err, "thresholds breached")
	})

	t.Run("No error if thresholds breached", func(t *testing.T) {
		buildCommand = MockCommandBuilderWithError(&MockThresholdError{})
		setFilePathEnvs()
		os.Setenv("PARAMETER_FAIL_ON_THRESHOLD_BREACH", "false")
		defer os.Clearenv()
		cfg, err := ConfigFromEnv()
		assert.NoError(t, err)
		err = RunPerfTests(cfg)
		assert.NoError(t, err)
	})
}

func TestReadLinesFromPipe(t *testing.T) {
	t.Run("Reads from pipe and closes", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
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
