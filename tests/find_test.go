package tests

import (
	"errors"
	"testing"

	"github.com/Silvelle/go-shell/src/shell"
	"github.com/Silvelle/go-shell/src/vfs"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{
			name:  "everything from the current directory",
			lines: []string{"find"},
			want:  ".\n./docs\n./docs/a.txt\n./docs/sub\n./docs/sub/b.md\n./empty\n./readme.txt",
		},
		{name: "start path", lines: []string{"find docs/sub"}, want: "docs/sub\ndocs/sub/b.md"},
		{name: "absolute start path", lines: []string{"cd docs", "find /docs/sub"}, want: "/docs/sub\n/docs/sub/b.md"},
		{name: "file as start path", lines: []string{"find readme.txt"}, want: "readme.txt"},
		{name: "name pattern", lines: []string{"find / -name *.txt"}, want: "/docs/a.txt\n/readme.txt"},
		{name: "name pattern without path", lines: []string{"find -name sub"}, want: "./docs/sub"},
		{name: "nothing matches", lines: []string{"find -name *.go"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustRun(t, newTestShell(), tt.lines...); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFindErrors(t *testing.T) {
	tests := []struct {
		input   string
		wantErr error
	}{
		{input: "find nope", wantErr: vfs.ErrNotFound},
		{input: "find -name", wantErr: shell.ErrInvalidArgs},
		{input: "find -size 10", wantErr: shell.ErrInvalidArgs},
		{input: "find . -name [", wantErr: shell.ErrInvalidArgs},
		{input: "find a b", wantErr: shell.ErrInvalidArgs},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if _, err := newTestShell().Execute(tt.input); !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
