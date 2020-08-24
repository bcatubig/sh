package sh

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// Args set the arguments passed to the binary
func Args(args ...string) func(*Command) {
	return func(c *Command) {
		c.args = args
	}
}

// Environment set the command subshell environment
func Environment(env map[string]string) func(*Command) {
	return func(c *Command) {
		c.environment = env
	}
}

// WorkingDir set the command subshell current working directory
func WorkingDir(dir string) func(*Command) {
	return func(c *Command) {
		c.workingDir = dir
	}
}

// Writers sets the writers used for stdout/stderr output
func Writers(writers []io.Writer) func(*Command) {
	return func(c *Command) {
		c.writers = writers
	}
}

// ExpectedReturnCode sets the expected process return code
//
// Default: 0
func ExpectedReturnCode(code int) func(*Command) {
	return func(c *Command) {
		c.expectedReturnCode = code
	}
}

// Command represents a command to be run
type Command struct {
	binary             string
	args               []string
	environment        map[string]string
	expectedReturnCode int
	workingDir         string
	writers            []io.Writer
}

// NewCommand generates a new Command object
func NewCommand(binary string, opts ...func(command *Command)) *Command {
	c := &Command{
		binary:             binary,
		args:               []string{},
		environment:        nil,
		expectedReturnCode: 0,
		workingDir:         "",
		writers:            []io.Writer{},
	}

	// Apply optional arguments
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Run will execute a shell command and wait for it to finish, returning stdout/stderr combined output
func (c *Command) Run() (*RunOutput, error) {

	// Create a buffer to write stdout/stderr to
	buf := &bytes.Buffer{}

	// Create a multiwriter
	mw := io.MultiWriter(append(c.writers, buf)...)

	// Create a low level Command object
	cmd := exec.Command(c.binary, c.args...)

	// Configure the command to write to our multi-writer
	cmd.Stdout = mw
	cmd.Stderr = mw

	// Configure the working directory
	cmd.Dir = c.workingDir

	// Configure Environment
	if c.environment != nil {
		envSlice := generateEnvSlice(c.environment)
		cmd.Env = envSlice
	}

	err := cmd.Run()

	if err != nil {
		// TODO: check for non-zero return code
		return &RunOutput{
			ReturnCode: 0,
			Output:     buf,
		}, err
	}

	return &RunOutput{
		ReturnCode: 0, // TODO: change this
		Output:     buf,
	}, nil

}

// RunOutput returns the process return code and combined stdout/stderr output
type RunOutput struct {
	ReturnCode int
	Output     *bytes.Buffer
}

func generateEnvSlice(env map[string]string) []string {
	var result []string
	for k, v := range env {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}
