package tests

import (
	"errors"
	"strings"
	"testing"

	"github.com/Silvelle/go-shell/src/shell"
)

// runScript runs a script with the given text and returns everything printed.
func runScript(t *testing.T, sh *shell.Shell, text string) (string, error) {
	t.Helper()
	path := writeFile(t, "start.txt", text)
	var lines []string
	err := sh.RunScript(path, func(s string) { lines = append(lines, s) })
	return strings.Join(lines, "\n"), err
}

func TestRunScriptShowsInputAndOutput(t *testing.T) {
	out, err := runScript(t, shell.New(), "# comment\n\nls /home\ncd /home\n")
	if err != nil {
		t.Fatalf("RunScript() error = %v", err)
	}
	for _, want := range []string{"$ ls /home", "$ cd /home"} {
		if !strings.Contains(out, want) {
			t.Errorf("output %q does not contain %q", out, want)
		}
	}
	if strings.Contains(out, "comment") {
		t.Errorf("comment lines must not be shown: %q", out)
	}
}

func TestRunScriptSkipsFailingLines(t *testing.T) {
	out, err := runScript(t, shell.New(), "foo\ncd a b\nls /home\n")
	if err != nil {
		t.Fatalf("RunScript() error = %v", err)
	}
	wantSkipped := 2
	if strings.Count(out, "line skipped") != wantSkipped {
		t.Errorf("want 2 reported errors, got output:\n%s", out)
	}
	if !strings.Contains(out, ":1:") || !strings.Contains(out, ":2:") {
		t.Errorf("errors must mention line numbers:\n%s", out)
	}
	if !strings.Contains(out, "$ ls /home") {
		t.Errorf("script must continue after errors:\n%s", out)
	}
}

func TestRunScriptStopsAtExit(t *testing.T) {
	out, err := runScript(t, shell.New(), "ls\nexit\nls never\n")
	if !errors.Is(err, shell.ErrExit) {
		t.Fatalf("RunScript() error = %v, want ErrExit", err)
	}
	if strings.Contains(out, "never") {
		t.Errorf("lines after exit must not run:\n%s", out)
	}
}

func TestRunScriptMissingFile(t *testing.T) {
	err := shell.New().RunScript("/no/such/script.txt", func(string) {})
	if err == nil {
		t.Error("RunScript() error = nil, want an error")
	}
}
