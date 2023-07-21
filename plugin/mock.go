// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package plugin

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type MockCommand struct {
	args    []string
	waitErr error
}

func (m *MockCommand) Start() error {
	return nil
}

func (m *MockCommand) Wait() error {
	return m.waitErr
}

func (m *MockCommand) String() (str string) {
	for _, arg := range m.args {
		str = fmt.Sprintf("%s %s", str, arg)
	}
	return
}

func (m *MockCommand) StdoutPipe() (io.ReadCloser, error) {
	dummyReader := strings.NewReader("")
	return io.NopCloser(dummyReader), nil
}

func (m *MockCommand) StderrPipe() (io.ReadCloser, error) {
	dummyReader := strings.NewReader("")
	return io.NopCloser(dummyReader), nil
}

// MockCommandBuilderWithError returns a function that will return a MockCommand
// which will return the specified waitErr on cmd.Wait().
func MockCommandBuilderWithError(waitErr error) func(string, ...string) shellCommand {
	return func(name string, args ...string) shellCommand {
		return &MockCommand{
			args:    append([]string{name}, args...),
			waitErr: waitErr,
		}
	}
}

type MockThresholdError struct {
	exec.ExitError
}

func (m *MockThresholdError) ExitCode() int {
	return THRESHOLDS_BREACHED_EXIT_CODE
}

func (m *MockThresholdError) Error() string {
	return "This is a mock threshold breach error"
}
