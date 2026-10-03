package tests

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Silvelle/go-shell/src/config"
)

// writeFile creates a file with the given text in a temporary directory.
func writeFile(t *testing.T, name, text string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadFlagsOnly(t *testing.T) {
	cfg, err := config.Load([]string{"--vfs", "v", "--log", "l.csv", "--script", "s.txt"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	want := config.Config{VFSPath: "v", LogPath: "l.csv", ScriptPath: "s.txt"}
	if cfg != want {
		t.Errorf("Load() = %+v, want %+v", cfg, want)
	}
}

func TestLoadConfigFile(t *testing.T) {
	ini := writeFile(t, "c.ini", "[emulator]\nvfs = /data\nlog = file.csv\nscript = start.txt\n")
	cfg, err := config.Load([]string{"--config", ini})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	want := config.Config{VFSPath: "/data", LogPath: "file.csv", ScriptPath: "start.txt", ConfigPath: ini}
	if cfg != want {
		t.Errorf("Load() = %+v, want %+v", cfg, want)
	}
}

func TestLoadFlagsOverrideConfigFile(t *testing.T) {
	ini := writeFile(t, "c.ini", "vfs = from-file\nlog = from-file.csv\n")
	cfg, err := config.Load([]string{"--config", ini, "--vfs", "from-flag"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.VFSPath != "from-flag" {
		t.Errorf("VFSPath = %q, want value from the flag", cfg.VFSPath)
	}
	if cfg.LogPath != "from-file.csv" {
		t.Errorf("LogPath = %q, want value from the file", cfg.LogPath)
	}
}

func TestLoadErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.ini")
	tests := []struct {
		name    string
		args    []string
		wantErr error
	}{
		{name: "missing config file", args: []string{"--config", missing}, wantErr: os.ErrNotExist},
		{
			name:    "bad INI line",
			args:    []string{"--config", writeFile(t, "a.ini", "oops\n")},
			wantErr: config.ErrINISyntax,
		},
		{
			name:    "unknown key",
			args:    []string{"--config", writeFile(t, "b.ini", "color = red\n")},
			wantErr: config.ErrUnknownKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := config.Load(tt.args)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Load() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadRejectsUnknownFlagAndExtraArgument(t *testing.T) {
	if _, err := config.Load([]string{"--color", "red"}); err == nil {
		t.Error("Load(--color) error = nil, want an error")
	}
	if _, err := config.Load([]string{"extra"}); err == nil {
		t.Error("Load(extra) error = nil, want an error")
	}
}

func TestParseINI(t *testing.T) {
	got, err := config.ParseINI("; comment\n# comment\n[section]\n\n  key = some value  \n")
	if err != nil {
		t.Fatalf("ParseINI() error = %v", err)
	}
	want := map[string]string{"key": "some value"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseINI() = %v, want %v", got, want)
	}
}

func TestDump(t *testing.T) {
	cfg := config.Config{VFSPath: "v", LogPath: "l", ScriptPath: "s", ConfigPath: "c"}
	want := []string{"vfs=v", "log=l", "script=s", "config=c"}
	if got := cfg.Dump(); !reflect.DeepEqual(got, want) {
		t.Errorf("Dump() = %v, want %v", got, want)
	}
}
