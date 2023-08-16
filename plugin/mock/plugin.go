package mock

import (
	"io"
	"os/exec"
	"strings"

	"github.com/go-vela/vela-k6/types"
)

const thresholdsBreachedExitCode = 99

type Command struct {
	args    []string
	waitErr error
}

func (m *Command) Start() error {
	return nil
}

func (m *Command) Wait() error {
	return m.waitErr
}

func (m *Command) String() (str string) {
	return ""
}

func (m *Command) StdoutPipe() (io.ReadCloser, error) {
	dummyReader := strings.NewReader("")
	return io.NopCloser(dummyReader), nil
}

func (m *Command) StderrPipe() (io.ReadCloser, error) {
	dummyReader := strings.NewReader("")
	return io.NopCloser(dummyReader), nil
}

// CommandBuilderWithError returns a function that will return a mock.Command
// which will return the specified waitErr on cmd.Wait().
func CommandBuilderWithError(waitErr error) func(string, ...string) types.ShellCommand {
	return func(name string, args ...string) types.ShellCommand {
		return &Command{
			args:    append([]string{name}, args...),
			waitErr: waitErr,
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
