package tests

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Silvelle/go-shell/src/vfs"
)

// makeTree creates files on disk under a temporary directory.
// Keys are slash-separated paths; a key ending in "/" is an empty directory.
func makeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if strings.HasSuffix(name, "/") {
			mustMkdir(t, path)
			continue
		}
		mustMkdir(t, filepath.Dir(path))
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// mustMkdir creates a directory with all its parents.
func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestVFSLoadReadsWholeTree(t *testing.T) {
	root := makeTree(t, map[string]string{
		"a.txt":         "hello",
		"docs/b.txt":    "world",
		"docs/x/y/z.md": "deep",
		"empty/":        "",
	})
	fs, err := vfs.Load(root)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	wantDirs, wantFiles := 4, 3
	dirs, files := fs.Count()
	if dirs != wantDirs || files != wantFiles {
		t.Errorf("Count() = %d dirs, %d files, want %d and %d", dirs, files, wantDirs, wantFiles)
	}
	node, err := fs.Resolve(fs.Root, "/docs/x/y/z.md")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if string(node.Content) != "deep" {
		t.Errorf("content = %q, want %q", node.Content, "deep")
	}
	if fs.Name != filepath.Base(root) {
		t.Errorf("Name = %q, want %q", fs.Name, filepath.Base(root))
	}
}

func TestVFSLoadErrors(t *testing.T) {
	file := makeTree(t, map[string]string{"f.txt": "x"})
	tests := []struct {
		name string
		path string
	}{
		{name: "path does not exist", path: filepath.Join(file, "missing")},
		{name: "path is a file, not a directory", path: filepath.Join(file, "f.txt")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := vfs.Load(tt.path); !errors.Is(err, vfs.ErrLoad) {
				t.Errorf("Load() error = %v, want ErrLoad", err)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	fs := vfs.Default()
	home, _ := fs.Resolve(fs.Root, "/home")
	tests := []struct {
		path string
		want string
	}{
		{path: "/", want: "/"},
		{path: "/home/user", want: "/home/user"},
		{path: "user", want: "/home/user"},
		{path: "./user/", want: "/home/user"},
		{path: "..", want: "/"},
		{path: "../..", want: "/"},
		{path: "user/../user//readme.txt", want: "/home/user/readme.txt"},
		{path: "~", want: "/"},
		{path: "~/home", want: "/home"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, err := fs.Resolve(home, tt.path)
			if err != nil {
				t.Fatalf("Resolve(%q) error = %v", tt.path, err)
			}
			if got.Path() != tt.want {
				t.Errorf("Resolve(%q) = %q, want %q", tt.path, got.Path(), tt.want)
			}
		})
	}
}

func TestResolveErrors(t *testing.T) {
	fs := vfs.Default()
	if _, err := fs.Resolve(fs.Root, "/nope"); !errors.Is(err, vfs.ErrNotFound) {
		t.Errorf("missing path: error = %v, want ErrNotFound", err)
	}
	if _, err := fs.Resolve(fs.Root, "/home/user/readme.txt/x"); !errors.Is(err, vfs.ErrNotDir) {
		t.Errorf("path through a file: error = %v, want ErrNotDir", err)
	}
}

func TestVFSLoadDoesNotModifyDisk(t *testing.T) {
	root := makeTree(t, map[string]string{"a.txt": "hello"})
	fs, err := vfs.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	fs.Root.Add(vfs.NewFile("new.txt", nil))
	if _, err := os.Stat(filepath.Join(root, "new.txt")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("changing the VFS must not touch the disk, stat error = %v", err)
	}
}