package tests

import (
	"errors"
	"strings"
	"testing"

	"github.com/Silvelle/go-shell/src/shell"
)

func TestExecuteStubs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "ls without arguments", input: "ls", want: "ls"},
		{name: "ls with arguments", input: "ls -l /home", want: "ls -l /home"},
		{name: "ls with extra spaces", input: "  ls   -a  ", want: "ls -a"},
		{name: "cd without arguments", input: "cd", want: "cd"},
		{name: "cd with one argument", input: "cd /tmp", want: "cd /tmp"},
		{name: "empty line does nothing", input: "   ", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shell.New().Execute(tt.input)
			if err != nil {
				t.Fatalf("Execute(%q) error = %v, want nil", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("Execute(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExecuteErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  error
		wantText string
	}{
		{name: "unknown command", input: "foo bar", wantErr: shell.ErrUnknownCommand, wantText: "foo"},
		{name: "cd with too many arguments", input: "cd a b", wantErr: shell.ErrInvalidArgs, wantText: "cd"},
		{name: "exit with arguments", input: "exit now", wantErr: shell.ErrInvalidArgs, wantText: "exit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := shell.New().Execute(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantText) {
				t.Errorf("Execute(%q) error = %q, want it to mention %q", tt.input, err, tt.wantText)
			}
		})
	}
}

func TestExecuteExit(t *testing.T) {
	out, err := shell.New().Execute("exit")
	if !errors.Is(err, shell.ErrExit) {
		t.Fatalf("Execute(\"exit\") error = %v, want ErrExit", err)
	}
	if out != "" {
		t.Errorf("Execute(\"exit\") output = %q, want empty", out)
	}
}
