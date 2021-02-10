package sh

import "io"

// CommandOption allows overriding various shell options
type CommandOption func(c *Command)

// WithArgs overrides arguments passed to the binary
func WithArgs(args ...string) func(*Command) {
	return func(c *Command) {
		c.args = args
	}
}

// WithEnvironment overrides the subshell environment
func WithEnvironment(env map[string]string) func(*Command) {
	return func(c *Command) {
		c.environment = env
	}
}

// WithStdIn overrides the default stdin stream to the command
func WithStdIn(input io.Reader) func(*Command) {
	return func(c *Command) {
		c.input = input
	}
}

// WithWorkingDir overrides the sub-shell's current working directory
func WithWorkingDir(dir string) func(*Command) {
	return func(c *Command) {
		c.workingDir = dir
	}
}

// WithWriters overrides IO writers used for stdout/stderr output
func WithWriters(writers ...io.Writer) func(*Command) {
	return func(c *Command) {
		c.writers = writers
	}
}

// WithExpectedReturnCode overrides the default expected process return code
//
// Default: 0
func WithExpectedReturnCode(code int) func(*Command) {
	return func(c *Command) {
		c.expectedReturnCode = code
	}
}

