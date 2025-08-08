// SPDX-License-Identifier: Apache-2.0

// Package mock provides mock implementations of the models.ShellCommand interface
// and a mock ThresholdError that simulates a threshold breach error.
package mock

import (
	"io"
	"os/exec"
	"strings"

	"github.com/go-vela/vela-k6/models"
)

const thresholdsBreachedExitCode = 99

// Command is a mock implementation of the models.ShellCommand interface.
type Command struct {
	args          []string
	waitErr       error
	stdoutPipeErr error
	stderrPipeErr error
	startErr      error
}

// Start is a mock implementation of the Start method.
func (m *Command) Start() error {
	return m.startErr
}

// Wait is a mock implementation of the Wait method.
func (m *Command) Wait() error {
	return m.waitErr
}

// String is a mock implementation of the String method.
func (m *Command) String() (str string) {
	return ""
}

// StdoutPipe is a mock implementation of the StdoutPipe method.
func (m *Command) StdoutPipe() (io.ReadCloser, error) {
	dummyReader := strings.NewReader("")
	return io.NopCloser(dummyReader), m.stdoutPipeErr
}

// StderrPipe is a mock implementation of the StderrPipe method.
func (m *Command) StderrPipe() (io.ReadCloser, error) {
	dummyReader := strings.NewReader("")
	return io.NopCloser(dummyReader), m.stderrPipeErr
}

// CommandBuilderWithError returns a function that will return a mock.Command
// which will return the specified waitErr on cmd.Wait().
func CommandBuilderWithError(waitErr error, stdoutPipeErr error, stderrPipeErr error, startErr error) func(string, ...string) models.ShellCommand {
	return func(name string, args ...string) models.ShellCommand {
		return &Command{
			args:          append([]string{name}, args...),
			waitErr:       waitErr,
			stdoutPipeErr: stdoutPipeErr,
			stderrPipeErr: stderrPipeErr,
			startErr:      startErr,
		}
	}
}

// ThresholdError is a mock implementation of the exec.ExitError interface
// that simulates a threshold breach error.
type ThresholdError struct {
	exec.ExitError
}

// ExitCode returns the exit code for the mock threshold breach error.
func (m *ThresholdError) ExitCode() int {
	return thresholdsBreachedExitCode
}

// Error returns a string representation of the mock threshold breach error.
func (m *ThresholdError) Error() string {
	return "This is a mock threshold breach error"
}
