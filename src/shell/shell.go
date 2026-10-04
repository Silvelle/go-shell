// Package shell implements the command interpreter of the emulator.
package shell

import "github.com/Silvelle/go-shell/src/vfs"

// handler runs one built-in command with its arguments.
type handler func(args []string) (string, error)

// EventLogger records every executed command.
type EventLogger interface {
	Log(command string, args []string, cmdErr error) error
}

// Shell holds the state of one emulator session.
type Shell struct {
	handlers map[string]handler
	logger   EventLogger
	fs       *vfs.VFS
	cwd      *vfs.Node
}

// New creates a Shell with all built-in commands registered,
// working in the default VFS.
func New() *Shell {
	s := &Shell{}
	s.handlers = map[string]handler{
		"ls":   s.ls,
		"cd":   s.cd,
		"exit": s.exit,
		"rev":  s.rev,
	}
	s.SetVFS(vfs.Default())
	return s
}

// SetLogger makes the shell record every command in l.
func (s *Shell) SetLogger(l EventLogger) {
	s.logger = l
}

// SetVFS switches the shell to fs and moves to its root.
func (s *Shell) SetVFS(fs *vfs.VFS) {
	s.fs = fs
	s.cwd = fs.Root
}

// Cwd returns the absolute path of the current directory.
func (s *Shell) Cwd() string {
	return s.cwd.Path()
}
