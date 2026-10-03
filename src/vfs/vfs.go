// Package vfs implements a virtual file system that lives in memory.
// It is loaded from a directory on disk, and the disk is never modified.
package vfs

import (
"errors"
"sort"
"strings"
)

// Errors returned by VFS operations. Their texts follow UNIX messages.
var (
	ErrNotFound = errors.New("no such file or directory")
	ErrNotDir   = errors.New("not a directory")
	ErrIsDir    = errors.New("is a directory")
	ErrLoad     = errors.New("cannot load VFS")
)

// Node is a file or a directory of the VFS.
type Node struct {
	Name     string
	IsDir    bool
	Content  []byte
	Parent   *Node
	Children map[string]*Node
}

// VFS is a tree of nodes kept entirely in memory.
type VFS struct {
	// Name is shown in the prompt, e.g. the source directory name.
	Name string
	Root *Node
}

// NewDir creates an empty directory node.
func NewDir(name string) *Node {
	return &Node{Name: name, IsDir: true, Children: map[string]*Node{}}
}

// NewFile creates a file node with the given content.
func NewFile(name string, content []byte) *Node {
	return &Node{Name: name, Content: content}
}

// Add puts child into the directory n.
func (n *Node) Add(child *Node) {
	child.Parent = n
	n.Children[child.Name] = child
}

// Path returns the absolute path of n inside the VFS.
func (n *Node) Path() string {
	if n.Parent == nil {
		return "/"
	}
	var parts []string
	for cur := n; cur.Parent != nil; cur = cur.Parent {
		parts = append([]string{cur.Name}, parts...)
	}
	return "/" + strings.Join(parts, "/")
}

// SortedChildren returns the children of n sorted by name.
func (n *Node) SortedChildren() []*Node {
	children := make([]*Node, 0, len(n.Children))
	for _, c := range n.Children {
		children = append(children, c)
	}
	sort.Slice(children, func(i, j int) bool { return children[i].Name < children[j].Name })
	return children
}

// Default returns a small VFS used when no directory is given.
func Default() *VFS {
	root := NewDir("")
	home := NewDir("home")
	user := NewDir("user")
	root.Add(home)
	home.Add(user)
	user.Add(NewFile("readme.txt", []byte("Welcome to the emulator!\n")))
	return &VFS{Name: "default", Root: root}
}

// Walk calls visit for n and every node below it, in name order.
// p is the path printed for n; children get p + "/" + name.
func Walk(n *Node, p string, visit func(p string, n *Node)) {
	visit(p, n)
	for _, c := range n.SortedChildren() {
		Walk(c, joinPath(p, c.Name), visit)
	}
}

// Count returns the number of directories (without the root) and files.
func (v *VFS) Count() (dirs, files int) {
	Walk(v.Root, "/", func(_ string, n *Node) {
		if n == v.Root {
			return
		}
		if n.IsDir {
			dirs++
		} else {
			files++
		}
	})
	return dirs, files
}

// joinPath appends name to p, adding a slash only when needed.
func joinPath(p, name string) string {
	if strings.HasSuffix(p, "/") {
		return p + name
	}
	return p + "/" + name
}