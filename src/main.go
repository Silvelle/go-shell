package main

import (
	"github.com/Silvelle/go-shell/src/shell"
	"github.com/Silvelle/go-shell/src/ui"
)

func main() {
	sh := shell.New()
	ui.Run(sh)
}
