// SPDX-License-Identifier: Apache-2.0

// Package models defines interfaces and types used in the Vela K6 plugin.
package models

import "io"

// ShellCommand is an interface that defines the methods for executing shell commands.
type ShellCommand interface {
	Start() error
	Wait() error
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
	String() string
}

// ErrorWithExitCode is an interface that defines a method for retrieving an exit code from an error.
type ErrorWithExitCode interface {
	ExitCode() int
}
