package vfs

import "errors"

// Errors returned by Remove and Move.
var (
	ErrRoot       = errors.New("cannot change the root directory")
	ErrExists     = errors.New("cannot overwrite: target exists and one of them is a directory")
	ErrIntoItself = errors.New("cannot move a directory into itself")
)

// Contains reports whether other is n itself or lies somewhere inside n.
func (n *Node) Contains(other *Node) bool {
	for cur := other; cur != nil; cur = cur.Parent {
		if cur == n {
			return true
		}
	}
	return false
}

// Remove deletes n, with everything inside it, from the VFS.
func (v *VFS) Remove(n *Node) error {
	if n.Parent == nil {
		return ErrRoot
	}
	delete(n.Parent.Children, n.Name)
	n.Parent = nil
	return nil
}

// Move puts n into the directory dir under the given name.
// An existing file with that name is replaced, like in UNIX mv;
// an existing directory is never replaced.
func (v *VFS) Move(n, dir *Node, name string) error {
	if n.Parent == nil {
		return ErrRoot
	}
	if n.Contains(dir) {
		return ErrIntoItself
	}
	old, exists := dir.Children[name]
	if exists && old == n {
		return nil
	}
	if exists && (old.IsDir || n.IsDir) {
		return ErrExists
	}

	delete(n.Parent.Children, n.Name)
	n.Name = name
	dir.Add(n)
	return nil
}