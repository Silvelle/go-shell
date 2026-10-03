// Command emulator starts the UNIX shell emulator with a graphical interface.
package main

import (
	"fmt"
	"os"

	"github.com/Silvelle/go-shell/src/config"
	"github.com/Silvelle/go-shell/src/logger"
	"github.com/Silvelle/go-shell/src/shell"
	"github.com/Silvelle/go-shell/src/ui"
)

// exitBadUsage is the exit code for invalid flags or config.
const exitBadUsage = 2

func main() {
	cfg, err := config.Load(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(exitBadUsage)
	}

	sh := shell.New()
	messages := settingsDump(cfg)
	if cfg.LogPath != "" {
		closeLog, err := attachLogger(sh, cfg.LogPath)
		if err != nil {
			messages = append(messages, "error: "+err.Error())
		} else {
			defer closeLog()
		}
	}
	ui.Run(sh, ui.Options{Messages: messages, ScriptPath: cfg.ScriptPath})
}

// settingsDump formats all settings for the debug output at startup.
func settingsDump(cfg config.Config) []string {
	lines := []string{"# settings"}
	for _, kv := range cfg.Dump() {
		lines = append(lines, "#   "+kv)
	}
	return lines
}

// attachLogger opens the CSV log and connects it to the shell.
// It returns a function that closes the log.
func attachLogger(sh *shell.Shell, path string) (func(), error) {
	lg, err := logger.Open(path, shell.CurrentUser())
	if err != nil {
		return nil, err
	}
	sh.SetLogger(lg)
	return func() { lg.Close() }, nil
}
