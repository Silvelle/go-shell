package shell

import (
	"errors"
	"fmt"
	"strings"
)

// ErrExit is returned by Execute when the user asks to leave the shell.
// The caller (for example, the GUI) decides how to shut down.
var ErrExit = errors.New("exit")

// ErrUnknownCommand is returned when the command name is not supported.
var ErrUnknownCommand = errors.New("unknown command")

// ErrInvalidArgs is returned when a command gets arguments it cannot accept.
var ErrInvalidArgs = errors.New("invalid arguments")

// Execute parses line and runs the command it contains.
// It returns the text to show to the user, or an error.
// Every command is recorded by the logger, if one is set.
func (s *Shell) Execute(line string) (string, error) {
	cmd, err := Parse(line)
	if err != nil {
		return "", err
	}
	if cmd.IsEmpty() {
		return "", nil
	}

	out, err := s.run(cmd)
	s.logEvent(cmd, err)
	return out, err
}

// run finds the handler for cmd and calls it.
func (s *Shell) run(cmd Command) (string, error) {
	run, ok := s.handlers[cmd.Name]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnknownCommand, cmd.Name)
	}
	return run(cmd.Args)
}

// logEvent records cmd and its error. Exiting is not an error.
// A failure to write the log must not break the shell, so it is ignored.
func (s *Shell) logEvent(cmd Command, err error) {
	if s.logger == nil {
		return
	}
	if errors.Is(err, ErrExit) {
		err = nil
	}
	_ = s.logger.Log(cmd.Name, cmd.Args, err)
}

// exit asks the caller to end the session. It takes no arguments.
func (s *Shell) exit(args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("%w: exit: unexpected argument %q", ErrInvalidArgs, args[0])
	}
	return "", ErrExit
}

// rejectOptions returns an error for the first argument that looks
// like an option ("-x"), for commands that support no options.
func rejectOptions(name string, args []string) error {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			return fmt.Errorf("%w: %s: invalid option '%s'", ErrInvalidArgs, name, a)
		}
	}
	return nil
}