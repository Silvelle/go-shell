package shell

import (
	"errors"
	"fmt"
	"strings"
)

// maxCdArgs is the largest number of arguments cd accepts.
const maxCdArgs = 1

// ErrExit is returned by Execute when the user asks  to leave the shell.
// The caller (for example, GUI) decides how to shut down.
var ErrExit = errors.New("exit")

// ErrUnknownCommand is returned when the command name is not supported.
var ErrUnknownCommand = errors.New("unknown command")

// ErrInvalidArgs is returned when a command gets arguments it cannot accept.
var ErrInvalidArgs = errors.New("invalid arguments")

// Execute parses line and runs the command it contains.
// It returns the text to show to the user, or an error.
func (s *Shell) Execute(line string) (string, error) {
	cmd, err := Parse(line)
	if err != nil {
		return "", err
	}
	if cmd.IsEmpty() {
		return "", nil
	}

	run, ok := s.handlers[cmd.Name]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnknownCommand, cmd.Name)
	}
	return run(cmd.Args)
}

// ls is a stub that echoes its name and arguments
func (s *Shell) ls(args []string) (string, error) {
	return stub("ls", args), nil
}

// cd is a stub that echoes its name and argument.
// Like a real shell, it rejects more than one argument.
func (s *Shell) cd(args []string) (string, error) {
	if len(args) > maxCdArgs {
		return "", fmt.Errorf("%w: cd: too many arguments", ErrInvalidArgs)
	}
	return stub("cd", args), nil
}

// exit asks the caller to end the session. It takes no arguments.
func (s *Shell) exit(args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("%w: exit: unexpexted argument %q", ErrInvalidArgs, args[0])
	}
	return "", ErrExit
}

// stub formats a command name and its arguments as one line.
func stub(name string, args []string) string {
	return strings.Join(append([]string{name}, args...), " ")
}
