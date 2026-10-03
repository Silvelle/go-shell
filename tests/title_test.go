package tests

import (
	"strings"
	"testing"

	"github.com/Silvelle/go-shell/src/shell"
)

func TestTitle(t *testing.T) {
	got := shell.Title("anna", "laptop")
	want := "Эмулятор - [anna@laptop]"
	if got != want {
		t.Errorf("Title() = %q, want %q", got, want)
	}
}

func TestSystemTitleHasUserAndHost(t *testing.T) {
	got := shell.SystemTitle()
	if !strings.HasPrefix(got, "Эмулятор - [") || !strings.Contains(got, "@") {
		t.Errorf("SystemTitle() = %q, want format \"Эмулятор - [user@host]\"", got)
	}
}
