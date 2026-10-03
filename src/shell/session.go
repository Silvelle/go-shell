package shell

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Prompt returns the text shown before every command.
func (s *Shell) Prompt() string {
	return "$ "
}

// Interact echoes line after the prompt, runs it and prints its output
// through out. The command's error, if any, is returned for the caller
// to report.
func (s *Shell) Interact(line string, out func(string)) error {
	out(s.Prompt() + line)
	result, err := s.Execute(line)
	if result != "" {
		out(result)
	}
	return err
}

// RunScript runs the startup script at path line by line, showing every
// command and its output through out, as if the user typed them.
// A failing line is reported and skipped. Blank lines and lines starting
// with '#' are ignored. It returns ErrExit if the script runs exit.
func (s *Shell) RunScript(path string, out func(string)) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("script: %w", err)
	}
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimRight(line, "\r")
		if isScriptComment(line) {
			continue
		}
		err := s.Interact(line, out)
		if errors.Is(err, ErrExit) {
			return err
		}
		if err != nil {
			out(fmt.Sprintf("error: %s:%d: %v (line skipped)", path, i+1, err))
		}
	}
	return nil
}

// isScriptComment reports whether a script line should not be run.
func isScriptComment(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed == "" || strings.HasPrefix(trimmed, "#")
}
