package shell

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/Silvelle/go-shell/src/vfs"
)

// rev prints every line of the given files with its characters reversed.
// Files that cannot be read are reported, the others are still printed.
func (s *Shell) rev(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("%w: rev: missing file operand", ErrInvalidArgs)
	}

	var parts []string
	var errs []error
	for _, p := range args {
		content, err := s.readFile(p)
		if err != nil {
			errs = append(errs, fmt.Errorf("rev: %s: %w", p, err))
			continue
		}
		parts = append(parts, reverseLines(string(content)))
	}
	return strings.Join(parts, "\n"), errors.Join(errs...)
}

// readFile returns the content of the file at path p.
func (s *Shell) readFile(p string) ([]byte, error) {
	node, err := s.fs.Resolve(s.cwd, p)
	if err != nil {
		return nil, err
	}
	if node.IsDir {
		return nil, vfs.ErrIsDir
	}
	return node.Content, nil
}

// reverseLines reverses the characters of every line of text.
func reverseLines(text string) string {
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	for i, line := range lines {
		lines[i] = reverse(line)
	}
	return strings.Join(lines, "\n")
}

// reverse reverses a string by characters (runes), not by bytes,
// so that non-ASCII text such as Cyrillic stays readable.
func reverse(s string) string {
	runes := []rune(s)
	slices.Reverse(runes)
	return string(runes)
}
