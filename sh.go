package sh

import (
	"bytes"
	"io"
)

// RunOptions allows configuring additional shell command options
type RunOptions struct {
	Environment map[string]string
	WorkingDir  string
	Writers     []io.Writer
}

// RunOutput returns the process return code and combined stdout/stderr output
type RunOutput struct {
	ReturnCode int
	Output     *bytes.Buffer
}

// Run will execute a shell command
func Run(binary string, args []string, options *RunOptions) (*RunOutput, error) {
	return nil, nil
}
