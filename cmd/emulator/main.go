package main

import (
	"github.com/Silvelle/go-shell/internal/shell"
	"github.com/Silvelle/go-shell/internal/ui"
)

func main() {
	sh := shell.New()
	ui.Run(sh)
}
