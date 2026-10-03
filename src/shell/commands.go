// Package shell implements the command interpreter of the emulator.
package shell

// handler runs one built-in command with its arguments
type handler func(args []string) (string, error)

// Shell holds the state of one emulator session.
type Shell struct {
	handlers map[string]handler
}

// New creates a Shell with all built-in commands registered
func New() *Shell {
	s := &Shell{}
	s.handlers = map[string]handler{
		"ls":   s.ls,
		"cd":   s.cd,
		"exit": s.exit,
	}
	return s
}
