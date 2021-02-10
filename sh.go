package sh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

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
	return runCommandWithContext(context.TODO(), c)
}

// Run will execute a shell command with a context and wait for it to finish, returning stdout/stderr combined output
func (c *Command) RunWithContext(ctx context.Context) (*RunOutput, error) {
	return runCommandWithContext(ctx, c)
}

func runCommandWithContext(ctx context.Context, c *Command) (*RunOutput, error) {
	// Create a buffer to write stdout/stderr to
	buf := &bytes.Buffer{}

	// Create a multi-writer
	mw := io.MultiWriter(append(c.writers, buf)...)

	// Create a low level Command object
	cmd := exec.CommandContext(ctx, c.binary, c.args...)

	// Configure the command to write to our multi-writer
	cmd.Stdout = mw
	cmd.Stderr = mw

	// Configure stdin
	cmd.Stdin = c.input

	// Configure the working directory
	cmd.Dir = c.workingDir

	// Configure the environment
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

	return -1, fmt.Errorf("not a valid ExitError")
}
