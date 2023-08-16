package types

import "io"

type ShellCommand interface {
	Start() error
	Wait() error
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
	String() string
}

type ErrorWithExitCode interface {
	ExitCode() int
}
