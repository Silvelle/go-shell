package vfs

import (
"fmt"
"os"
"path/filepath"
)

// Load reads the directory at path and everything inside it into memory.
// It fails if path does not exist or is not a directory.
func Load(path string) (*VFS, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLoad, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%w: %s: %v (VFS must be a directory)", ErrLoad, path, ErrNotDir)
	}

	root := NewDir("")
	if err := loadDir(path, root); err != nil {
		return nil, err
	}
	return &VFS{Name: filepath.Base(filepath.Clean(path)), Root: root}, nil
}

// loadDir adds every entry of the disk directory dir to node.
func loadDir(dir string, node *Node) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrLoad, err)
	}
	for _, e := range entries {
		child, err := loadEntry(filepath.Join(dir, e.Name()), e)
		if err != nil {
			return err
		}
		node.Add(child)
	}
	return nil
}

// loadEntry reads one file or, recursively, one directory.
func loadEntry(path string, e os.DirEntry) (*Node, error) {
	if e.IsDir() {
		child := NewDir(e.Name())
		return child, loadDir(path, child)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLoad, err)
	}
	return NewFile(e.Name(), data), nil
}