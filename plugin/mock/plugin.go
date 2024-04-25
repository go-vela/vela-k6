package mock

import (
	"io"
	"os/exec"
	"strings"

	"github.com/go-vela/vela-k6/types"
)

const thresholdsBreachedExitCode = 99

type Command struct {
	args          []string
	waitErr       error
	stdoutPipeErr error
	stderrPipeErr error
	startErr      error
}

func (m *Command) Start() error {
	return m.startErr
}

func (m *Command) Wait() error {
	return m.waitErr
}

func (m *Command) String() (str string) {
	return ""
}

func (m *Command) StdoutPipe() (io.ReadCloser, error) {
	dummyReader := strings.NewReader("")
	return io.NopCloser(dummyReader), m.stdoutPipeErr
}

func (m *Command) StderrPipe() (io.ReadCloser, error) {
	dummyReader := strings.NewReader("")
	return io.NopCloser(dummyReader), m.stderrPipeErr
}

// CommandBuilderWithError returns a function that will return a mock.Command
// which will return the specified waitErr on cmd.Wait().
func CommandBuilderWithError(waitErr error, stdoutPipeErr error, stderrPipeErr error, startErr error) func(string, ...string) types.ShellCommand {
	return func(name string, args ...string) types.ShellCommand {
		return &Command{
			args:          append([]string{name}, args...),
			waitErr:       waitErr,
			stdoutPipeErr: stdoutPipeErr,
			stderrPipeErr: stderrPipeErr,
			startErr:      startErr,
		}
	}
}

type ThresholdError struct {
	exec.ExitError
}

func (m *ThresholdError) ExitCode() int {
	return thresholdsBreachedExitCode
}

func (m *ThresholdError) Error() string {
	return "This is a mock threshold breach error"
}
