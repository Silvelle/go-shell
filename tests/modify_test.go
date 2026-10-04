package tests

import (
	"errors"
	"testing"

	"github.com/Silvelle/go-shell/src/shell"
	"github.com/Silvelle/go-shell/src/vfs"
)

func TestRm(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{name: "file", lines: []string{"rm readme.txt", "ls"}, want: "docs/  empty/"},
		{name: "several files", lines: []string{"rm docs/a.txt docs/sub/b.md", "find docs"}, want: "docs\ndocs/sub"},
		{name: "directory with -r", lines: []string{"rm -r docs", "ls"}, want: "empty/  readme.txt"},
		{name: "directory with -R", lines: []string{"rm -R empty", "ls"}, want: "docs/  readme.txt"},
		{name: "relative from subdirectory", lines: []string{"cd docs", "rm ../readme.txt", "ls /"}, want: "docs/  empty/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustRun(t, newTestShell(), tt.lines...); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRmErrors(t *testing.T) {
	tests := []struct {
		name    string
		lines   []string
		input   string
		wantErr error
	}{
		{name: "no operand", input: "rm", wantErr: shell.ErrInvalidArgs},
		{name: "only flag", input: "rm -r", wantErr: shell.ErrInvalidArgs},
		{name: "missing path", input: "rm nope", wantErr: vfs.ErrNotFound},
		{name: "directory without -r", input: "rm docs", wantErr: vfs.ErrIsDir},
		{name: "root", input: "rm -r /", wantErr: vfs.ErrRoot},
		{name: "current directory", lines: []string{"cd docs/sub"}, input: "rm -r /docs", wantErr: shell.ErrBusy},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh := newTestShell()
			mustRun(t, sh, tt.lines...)
			if _, err := sh.Execute(tt.input); !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestRmRemovesOthersDespiteErrors(t *testing.T) {
	sh := newTestShell()
	if _, err := sh.Execute("rm nope readme.txt"); err == nil {
		t.Error("error = nil, want an error for 'nope'")
	}
	if got := mustRun(t, sh, "ls"); got != "docs/  empty/" {
		t.Errorf("readme.txt must be removed, ls = %q", got)
	}
}

func TestMv(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{name: "rename file", lines: []string{"mv readme.txt hello.txt", "ls"}, want: "docs/  empty/  hello.txt"},
		{name: "file into directory", lines: []string{"mv readme.txt empty", "ls empty"}, want: "readme.txt"},
		{name: "file into directory with new name", lines: []string{"mv readme.txt empty/r.txt", "ls empty"}, want: "r.txt"},
		{
			name:  "directory into directory",
			lines: []string{"mv docs empty", "find empty"},
			want:  "empty\nempty/docs\nempty/docs/a.txt\nempty/docs/sub\nempty/docs/sub/b.md",
		},
		{name: "rename directory", lines: []string{"mv docs papers", "ls"}, want: "empty/  papers/  readme.txt"},
		{name: "replace existing file", lines: []string{"mv readme.txt docs/a.txt", "rev docs/a.txt"}, want: "тевирП"},
		{name: "to parent with dot dot", lines: []string{"cd docs/sub", "mv b.md ..", "ls /docs"}, want: "a.txt  b.md  sub/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustRun(t, newTestShell(), tt.lines...); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMvKeepsCurrentDirectoryValid(t *testing.T) {
	sh := newTestShell()
	mustRun(t, sh, "cd docs/sub", "mv /docs /papers")
	if sh.Cwd() != "/papers/sub" {
		t.Errorf("Cwd() = %q, want /papers/sub", sh.Cwd())
	}
}

func TestMvErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{name: "one argument", input: "mv readme.txt", wantErr: shell.ErrInvalidArgs},
		{name: "three arguments", input: "mv a b c", wantErr: shell.ErrInvalidArgs},
		{name: "missing source", input: "mv nope x", wantErr: vfs.ErrNotFound},
		{name: "missing destination directory", input: "mv readme.txt nope/x", wantErr: vfs.ErrNotFound},
		{name: "destination parent is a file", input: "mv docs/a.txt readme.txt/x", wantErr: vfs.ErrNotDir},
		{name: "directory into itself", input: "mv docs docs/sub", wantErr: vfs.ErrIntoItself},
		{name: "directory over a file", input: "mv empty readme.txt", wantErr: vfs.ErrExists},
		{name: "file over a directory", input: "mv docs/sub/b.md docs/sub", wantErr: nil},
		{name: "root", input: "mv / x", wantErr: vfs.ErrRoot},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := newTestShell().Execute(tt.input); !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestChangesDoNotTouchDisk(t *testing.T) {
	root := makeTree(t, map[string]string{"a.txt": "x", "d/b.txt": "y"})
	fs, err := vfs.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	sh := shell.New()
	sh.SetVFS(fs)
	mustRun(t, sh, "rm a.txt", "mv d e")

	again, err := vfs.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	wantDirs, wantFiles := 1, 2
	if dirs, files := again.Count(); dirs != wantDirs || files != wantFiles {
		t.Errorf("disk changed: %d dirs, %d files, want %d and %d", dirs, files, wantDirs, wantFiles)
	}
}