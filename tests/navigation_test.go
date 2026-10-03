package tests

import (
	"errors"
	"strings"
	"testing"

	"github.com/Silvelle/go-shell/src/shell"
	"github.com/Silvelle/go-shell/src/vfs"
)


func TestLs(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{name: "current directory", lines: []string{"ls"}, want: "docs/  empty/  readme.txt"},
		{name: "absolute path", lines: []string{"ls /docs"}, want: "a.txt  sub/"},
		{name: "relative path after cd", lines: []string{"cd docs", "ls sub"}, want: "b.md"},
		{name: "parent directory", lines: []string{"cd /docs/sub", "ls .."}, want: "a.txt  sub/"},
		{name: "empty directory", lines: []string{"ls empty"}, want: ""},
		{name: "file", lines: []string{"ls docs/a.txt"}, want: "docs/a.txt"},
		{name: "several paths", lines: []string{"ls docs/sub empty"}, want: "docs/sub:\nb.md\n\nempty:\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustRun(t, newTestShell(), tt.lines...); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLsErrors(t *testing.T) {
	sh := newTestShell()
	out, err := sh.Execute("ls nope empty docs/sub")
	if !errors.Is(err, vfs.ErrNotFound) || !strings.Contains(err.Error(), "'nope'") {
		t.Errorf("error = %v, want ErrNotFound mentioning 'nope'", err)
	}
	if !strings.Contains(out, "docs/sub:\nb.md") {
		t.Errorf("existing paths must still be listed, got %q", out)
	}
	if _, err := sh.Execute("ls -l"); !errors.Is(err, shell.ErrInvalidArgs) {
		t.Errorf("ls -l error = %v, want ErrInvalidArgs", err)
	}
}

func TestCd(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{name: "absolute", lines: []string{"cd /docs/sub"}, want: "/docs/sub"},
		{name: "relative, step by step", lines: []string{"cd docs", "cd sub"}, want: "/docs/sub"},
		{name: "dot dot", lines: []string{"cd /docs/sub", "cd .."}, want: "/docs"},
		{name: "dot dot at the root stays", lines: []string{"cd ..", "cd ../.."}, want: "/"},
		{name: "no arguments goes home", lines: []string{"cd /docs", "cd"}, want: "/"},
		{name: "tilde", lines: []string{"cd /docs", "cd ~/empty"}, want: "/empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh := newTestShell()
			mustRun(t, sh, tt.lines...)
			if sh.Cwd() != tt.want {
				t.Errorf("Cwd() = %q, want %q", sh.Cwd(), tt.want)
			}
		})
	}
}

func TestCdErrors(t *testing.T) {
	tests := []struct {
		input   string
		wantErr error
	}{
		{input: "cd nope", wantErr: vfs.ErrNotFound},
		{input: "cd ../readme.txt", wantErr: vfs.ErrNotDir},
		{input: "cd a.txt/x", wantErr: vfs.ErrNotDir},
		{input: "cd docs empty", wantErr: shell.ErrInvalidArgs},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			sh := newTestShell()
			mustRun(t, sh, "cd docs")
			if _, err := sh.Execute(tt.input); !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
			if sh.Cwd() != "/docs" {
				t.Errorf("a failed cd must not move, Cwd() = %q", sh.Cwd())
			}
		})
	}
}