// Command emulator starts the UNIX shell emulator with a graphical interface.
package main

import (
	"github.com/Silvelle/go-shell/src/shell"
	"github.com/Silvelle/go-shell/src/ui"
)

func main() {
	ui.Run(shell.New())
}
