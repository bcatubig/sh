package sh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

// Args sets the arguments passed to the binary
func Args(args ...string) func(*Command) {
	return func(c *Command) {
		c.args = args
	}
}

// Environment sets the command subshell environment
func Environment(env map[string]string) func(*Command) {
	return func(c *Command) {
		c.environment = env
	}
}

// Input sets stdin input to the process
func Input(input io.Reader) func(*Command) {
	return func(c *Command) {
		c.input = input
	}
}

// WorkingDir sets the command subshell current working directory
func WorkingDir(dir string) func(*Command) {
	return func(c *Command) {
		c.workingDir = dir
	}
}

// Writers sets the writers used for stdout/stderr output
func Writers(writers ...io.Writer) func(*Command) {
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
	input              io.Reader
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
		input:              nil,
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
	return runCommand(nil, c)
}

// Run will execute a shell command with a context and wait for it to finish, returning stdout/stderr combined output
func (c *Command) RunWithContext(ctx context.Context) (*RunOutput, error) {
	return runCommand(ctx, c)
}

func runCommand(ctx context.Context, c *Command) (*RunOutput, error) {
	// Create a buffer to write stdout/stderr to
	buf := &bytes.Buffer{}

	// Create a multiwriter
	mw := io.MultiWriter(append(c.writers, buf)...)

	// Create a low level Command object
	var cmd *exec.Cmd
	if ctx == nil {
		cmd = exec.Command(c.binary, c.args...)
	} else {
		cmd = exec.CommandContext(ctx, c.binary, c.args...)
	}

	// Configure the command to write to our multi-writer
	cmd.Stdout = mw
	cmd.Stderr = mw

	// Configure stdin
	cmd.Stdin = c.input

	// Configure the working directory
	cmd.Dir = c.workingDir

	// Configure Environment
	if c.environment != nil {
		envSlice := generateEnvSlice(c.environment)
		cmd.Env = envSlice
	}

	err := cmd.Run()

	if err != nil {
		// Validate that non-zero is equal to expectedExitCode if set
		rc, rcErr := getReturnCode(err)

		if rcErr != nil {
			return &RunOutput{
				ReturnCode: 1,
				Output:     buf,
			}, err
		}

		if rc != c.expectedReturnCode {
			return &RunOutput{
				ReturnCode: rc,
				Output:     buf,
			}, fmt.Errorf("unexpected exit code: expected: %d, got %d", c.expectedReturnCode, rc)
		}

		return &RunOutput{
			ReturnCode: rc,
			Output:     buf,
		}, nil
	}

	return &RunOutput{
		ReturnCode: 0,
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

func getReturnCode(err error) (int, error) {
	var e *exec.ExitError
	if errors.As(err, &e) {
		return e.ExitCode(), nil
	}

	return 0, fmt.Errorf("not a valid ExitError")
}
