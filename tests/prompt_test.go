package tests

import (
	"testing"

	"github.com/Silvelle/go-shell/src/shell"
	"github.com/Silvelle/go-shell/src/vfs"
)

func TestPromptShowsVFSNameAndCwd(t *testing.T) {
	sh := shell.New()
	if got, want := sh.Prompt(), "default:/$ "; got != want {
		t.Errorf("Prompt() = %q, want %q", got, want)
	}

	root := vfs.NewDir("")
	sh.SetVFS(&vfs.VFS{Name: "mine", Root: root})
	if got, want := sh.Prompt(), "mine:/$ "; got != want {
		t.Errorf("Prompt() after SetVFS = %q, want %q", got, want)
	}
}