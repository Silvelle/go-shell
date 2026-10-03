// Package config reads emulator settings from command line flags
// and an INI configuration file.
package config

import (
	"flag"
	"fmt"
)

// Settings names, used both as flag names and as INI keys.
const (
	KeyVFS    = "vfs"
	KeyLog    = "log"
	KeyScript = "script"
	KeyConfig = "config"
)

// Config hold all emulator settings.
type Config struct {
	VFSPath    string
	LogPath    string
	ScriptPath string
	ConfigPath string
}

// Load build a Config from command-line arguments (without the program name).
// If a config file is given, it is read too, but values from the command line
// take priority over values from the file.
func Load(args []string) (Config, error) {
	var cfg Config
	fs := flag.NewFlagSet("emulator", flag.ContinueOnError)
	fs.StringVar(&cfg.VFSPath, KeyVFS, "", "path to the VFS directory")
	fs.StringVar(&cfg.LogPath, KeyLog, "", "path to the CSV log file")
	fs.StringVar(&cfg.ScriptPath, KeyScript, "", "path to the startup script")
	fs.StringVar(&cfg.ConfigPath, KeyConfig, "", "path to the INI config file")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	if fs.NArg() != 0 {
		return Config{}, fmt.Errorf("unexpected argument %q", fs.Arg(0))
	}
	if cfg.ConfigPath == "" {
		return cfg, nil
	}

	values, err := ReadINI(cfg.ConfigPath)
	if err != nil {
		return Config{}, err
	}
	fromFlags := map[string]bool{}
	fs.Visit(func(f *flag.Flag) { fromFlags[f.Name] = true })
	return cfg, cfg.fill(values, fromFlags)
}

// fill copies values from the config file into settings
// that were not given on the command line.
func (c *Config) fill(values map[string]string, fromFlags map[string]bool) error {
	fields := map[string]*string{
		KeyVFS:    &c.VFSPath,
		KeyLog:    &c.LogPath,
		KeyScript: &c.ScriptPath,
	}
	for key, value := range values {
		field, ok := fields[key]
		if !ok {
			return fmt.Errorf("%w: %q in %s", ErrUnknownKey, key, c.ConfigPath)
		}
		if !fromFlags[key] {
			*field = value
		}
	}
	return nil
}

// Dump returns all settings as "key=value" lines in a fixed order.
func (c Config) Dump() []string {
	return []string{
		KeyVFS + "=" + c.VFSPath,
		KeyLog + "=" + c.LogPath,
		KeyScript + "=" + c.ScriptPath,
		KeyConfig + "=" + c.ConfigPath,
	}
}
