package vfs

import "strings"

// Resolve finds the node for path p. Relative paths start at cwd.
// It understands "/", ".", ".." and "~" (the root). ".." at the root
// stays at the root, like in UNIX.
func (v *VFS) Resolve(cwd *Node, p string) (*Node, error) {
	cur := cwd
	if p == "~" || strings.HasPrefix(p, "~/") {
		p = "/" + strings.TrimPrefix(p, "~")
	}
	if strings.HasPrefix(p, "/") {
		cur = v.Root
	}
	for _, part := range strings.Split(p, "/") {
		next, err := step(cur, part)
		if err != nil {
			return nil, err
		}
		cur = next
	}
	return cur, nil
}

// step moves from cur by one path component.
func step(cur *Node, part string) (*Node, error) {
	switch part {
	case "", ".":
		return cur, nil
	case "..":
		if cur.Parent == nil {
			return cur, nil
		}
		return cur.Parent, nil
	}
	if !cur.IsDir {
		return nil, ErrNotDir
	}
	child, ok := cur.Children[part]
	if !ok {
		return nil, ErrNotFound
	}
	return child, nil
}
