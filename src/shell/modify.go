package shell

import (
	"errors"
	"fmt"
	"path"

	"github.com/Silvelle/go-shell/src/vfs"
)

// mvArgs is the number of arguments mv takes: source and destination.
const mvArgs = 2

// ErrBusy is returned when removing a directory that holds the current one.
var ErrBusy = errors.New("the current directory is inside it")

// rm removes files. Directories are removed only with -r (or -R),
// together with everything inside them. Changes stay in memory.
func (s *Shell) rm(args []string) (string, error) {
	recursive, paths := splitRecursiveFlag(args)
	if len(paths) == 0 {
		return "", fmt.Errorf("%w: rm: missing operand", ErrInvalidArgs)
	}

	var errs []error
	for _, p := range paths {
		if err := s.remove(p, recursive); err != nil {
			errs = append(errs, fmt.Errorf("rm: cannot remove '%s': %w", p, err))
		}
	}
	return "", errors.Join(errs...)
}

// splitRecursiveFlag separates -r / -R from the paths.
func splitRecursiveFlag(args []string) (bool, []string) {
	recursive := false
	var paths []string
	for _, a := range args {
		if a == "-r" || a == "-R" {
			recursive = true
			continue
		}
		paths = append(paths, a)
	}
	return recursive, paths
}

// remove deletes the node at p after checking that it is allowed.
func (s *Shell) remove(p string, recursive bool) error {
	node, err := s.fs.Resolve(s.cwd, p)
	if err != nil {
		return err
	}
	if node.Parent == nil {
		return vfs.ErrRoot
	}
	if node.IsDir && !recursive {
		return vfs.ErrIsDir
	}
	if node.Contains(s.cwd) {
		return ErrBusy
	}
	return s.fs.Remove(node)
}

// mv moves or renames a file or a directory. If the destination is an
// existing directory, the source is moved into it. Changes stay in memory.
func (s *Shell) mv(args []string) (string, error) {
	if len(args) != mvArgs {
		return "", fmt.Errorf("%w: mv: usage: mv source destination", ErrInvalidArgs)
	}
	src, dst := args[0], args[1]

	node, err := s.fs.Resolve(s.cwd, src)
	if err != nil {
		return "", fmt.Errorf("mv: cannot stat '%s': %w", src, err)
	}
	dir, name, err := s.destination(dst, node.Name)
	if err == nil {
		err = s.fs.Move(node, dir, name)
	}
	if err != nil {
		return "", fmt.Errorf("mv: cannot move '%s' to '%s': %w", src, dst, err)
	}
	return "", nil
}

// destination works out the directory and the new name for mv.
// An existing directory dst keeps the old name; otherwise the last
// element of dst becomes the new name inside its parent directory.
func (s *Shell) destination(dst, name string) (*vfs.Node, string, error) {
	if node, err := s.fs.Resolve(s.cwd, dst); err == nil && node.IsDir {
		return node, name, nil
	}
	parentPath, base := path.Split(dst)
	if base == "" || base == "." || base == ".." {
		return nil, "", vfs.ErrNotFound
	}
	parent, err := s.fs.Resolve(s.cwd, parentPath)
	if err != nil {
		return nil, "", err
	}
	if !parent.IsDir {
		return nil, "", vfs.ErrNotDir
	}
	return parent, base, nil
}