package tests

import (
	"errors"
	"testing"

	"github.com/Silvelle/go-shell/src/shell"
	"github.com/Silvelle/go-shell/src/vfs"
)

func TestRev(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "rev docs/a.txt", want: "olleh\ndlrow"},
		{input: "rev readme.txt", want: "тевирП"},
		{input: "rev docs/sub/b.md readme.txt", want: "cba\nтевирП"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := mustRun(t, newTestShell(), tt.input); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRevErrors(t *testing.T) {
	tests := []struct {
		input   string
		wantErr error
	}{
		{input: "rev", wantErr: shell.ErrInvalidArgs},
		{input: "rev nope", wantErr: vfs.ErrNotFound},
		{input: "rev docs", wantErr: vfs.ErrIsDir},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if _, err := newTestShell().Execute(tt.input); !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestRevPrintsReadableFilesDespiteErrors(t *testing.T) {
	out, err := newTestShell().Execute("rev nope docs/sub/b.md")
	if err == nil || out != "cba" {
		t.Errorf("got %q, %v; want \"cba\" and an error", out, err)
	}
}
