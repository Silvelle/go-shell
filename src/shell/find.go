package shell

import (
	"fmt"
	"path"
	"strings"

	"github.com/Silvelle/go-shell/src/vfs"
)

// findNameArgs is the number of arguments in "-name pattern".
const findNameArgs = 2

// errFindUsage is returned for arguments find does not understand.
var errFindUsage = fmt.Errorf("%w: find: usage: find [path] [-name pattern]", ErrInvalidArgs)

// find prints the start path and every path below it, one per line.
// With -name, only paths whose last element matches the shell pattern
// (*, ?, [...]) are printed. The start path defaults to ".".
func (s *Shell) find(args []string) (string, error) {
	start, pattern, err := parseFindArgs(args)
	if err != nil {
		return "", err
	}
	node, err := s.fs.Resolve(s.cwd, start)
	if err != nil {
		return "", fmt.Errorf("find: '%s': %w", start, err)
	}

	var found []string
	vfs.Walk(node, start, func(p string, _ *vfs.Node) {
		if matches(pattern, p) {
			found = append(found, p)
		}
	})
	return strings.Join(found, "\n"), nil
}

// parseFindArgs splits find arguments into the start path and the pattern.
func parseFindArgs(args []string) (start, pattern string, err error) {
	start = "."
	if len(args) != 0 && !strings.HasPrefix(args[0], "-") {
		start, args = args[0], args[1:]
	}
	if len(args) == 0 {
		return start, "", nil
	}
	if len(args) != findNameArgs || args[0] != "-name" {
		return "", "", errFindUsage
	}
	if _, err := path.Match(args[1], ""); err != nil {
		return "", "", fmt.Errorf("%w: find: bad pattern %q", ErrInvalidArgs, args[1])
	}
	return start, args[1], nil
}

// matches reports whether the last element of p matches pattern.
// An empty pattern matches everything.
func matches(pattern, p string) bool {
	if pattern == "" {
		return true
	}
	ok, _ := path.Match(pattern, path.Base(p))
	return ok
}