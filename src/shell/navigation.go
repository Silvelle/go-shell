package shell

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Silvelle/go-shell/src/vfs"
)

// maxCdArgs is the largest number of arguments cd accepts.
const maxCdArgs = 1

// singlePath is the number of ls operands that are listed without headers.
const singlePath = 1

// ls lists directories and shows files. Without arguments it lists the
// current directory. With several paths every listing gets a header.
// Paths that cannot be found are reported, the others are still listed.
func (s *Shell) ls(args []string) (string, error) {
	if err := rejectOptions("ls", args); err != nil {
		return "", err
	}
	if len(args) == 0 {
		args = []string{"."}
	}

	var blocks []string
	var errs []error
	for _, p := range args {
		node, err := s.fs.Resolve(s.cwd, p)
		if err != nil {
			errs = append(errs, fmt.Errorf("ls: cannot access '%s': %w", p, err))
			continue
		}
		blocks = append(blocks, listing(p, node, len(args) > singlePath))
	}
	return strings.Join(blocks, "\n\n"), errors.Join(errs...)
}

// listing formats one ls operand: a file name, or the directory contents
// with an optional "path:" header.
func listing(p string, n *vfs.Node, header bool) string {
	if !n.IsDir {
		return p
	}
	names := make([]string, 0, len(n.Children))
	for _, c := range n.SortedChildren() {
		names = append(names, displayName(c))
	}
	body := strings.Join(names, "  ")
	if header {
		return p + ":\n" + body
	}
	return body
}

// displayName returns the name of n, with "/" added for directories.
func displayName(n *vfs.Node) string {
	if n.IsDir {
		return n.Name + "/"
	}
	return n.Name
}

// cd changes the current directory. Without arguments it goes to the
// home directory, which is the VFS root.
func (s *Shell) cd(args []string) (string, error) {
	if len(args) > maxCdArgs {
		return "", fmt.Errorf("%w: cd: too many arguments", ErrInvalidArgs)
	}
	target := "~"
	if len(args) == maxCdArgs {
		target = args[0]
	}

	node, err := s.fs.Resolve(s.cwd, target)
	if err != nil {
		return "", fmt.Errorf("cd: %s: %w", target, err)
	}
	if !node.IsDir {
		return "", fmt.Errorf("cd: %s: %w", target, vfs.ErrNotDir)
	}
	s.cwd = node
	return "", nil
}
