package tests

import (
	"testing"

	"github.com/Silvelle/go-shell/src/shell"
	"github.com/Silvelle/go-shell/src/vfs"
)

// newTestShell returns a shell working in this VFS:
func newTestShell() *shell.Shell {
	root := vfs.NewDir("")
	docs := vfs.NewDir("docs")
	sub := vfs.NewDir("sub")
	root.Add(docs)
	root.Add(vfs.NewDir("empty"))
	root.Add(vfs.NewFile("readme.txt", []byte("Привет\n")))
	docs.Add(vfs.NewFile("a.txt", []byte("hello\nworld\n")))
	docs.Add(sub)
	sub.Add(vfs.NewFile("b.md", []byte("abc")))

	sh := shell.New()
	sh.SetVFS(&vfs.VFS{Name: "test", Root: root})
	return sh
}

// mustRun runs lines one by one and fails the test on the first error.
// It returns the output of the last line.
func mustRun(t *testing.T, sh *shell.Shell, lines ...string) string {
	t.Helper()
	var out string
	for _, line := range lines {
		var err error
		out, err = sh.Execute(line)
		if err != nil {
			t.Fatalf("Execute(%q) error = %v", line, err)
		}
	}
	return out
}