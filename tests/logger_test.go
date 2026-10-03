package tests

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Silvelle/go-shell/src/logger"
	"github.com/Silvelle/go-shell/src/shell"
)

// wantRows is the header plus two logged commands.
const wantRows = 3

// readCSV returns all rows of a CSV file.
func readCSV(t *testing.T, path string) [][]string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

func TestLoggerWritesHeaderAndRows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.csv")
	lg, err := logger.Open(path, "anna")
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	lg.Log("ls", []string{"-l", "/home"}, nil)
	lg.Log("foo", nil, errors.New("unknown command: foo"))
	lg.Close()

	rows := readCSV(t, path)
	if len(rows) != wantRows {
		t.Fatalf("got %d rows, want header + 2", len(rows))
	}
	if rows[0][2] != "command" || rows[0][4] != "error" {
		t.Errorf("header = %v", rows[0])
	}
	if rows[1][1] != "anna" || rows[1][2] != "ls" || rows[1][3] != "-l /home" || rows[1][4] != "-" {
		t.Errorf("row 1 = %v", rows[1])
	}
	if rows[2][4] != "unknown command: foo" {
		t.Errorf("row 2 error = %q", rows[2][4])
	}
}

func TestLoggerAppendsWithoutSecondHeader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.csv")
	for range 2 {
		lg, err := logger.Open(path, "anna")
		if err != nil {
			t.Fatal(err)
		}
		lg.Log("ls", nil, nil)
		lg.Close()
	}
	if rows := readCSV(t, path); len(rows) != wantRows {
		t.Errorf("got %d rows, want header + 2", len(rows))
	}
}

func TestLoggerOpenError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no-such-dir", "log.csv")
	if _, err := logger.Open(path, "anna"); err == nil {
		t.Error("Open() error = nil, want an error")
	}
}

func TestShellLogsCommands(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.csv")
	lg, err := logger.Open(path, "anna")
	if err != nil {
		t.Fatal(err)
	}
	sh := shell.New()
	sh.SetLogger(lg)
	sh.Execute("ls /tmp")
	sh.Execute("foo")
	sh.Execute("   ")
	lg.Close()

	rows := readCSV(t, path)
	if len(rows) != wantRows {
		t.Fatalf("got %d rows, want header + 2 (blank lines are not logged)", len(rows))
	}
	if rows[2][2] != "foo" || rows[2][4] == "" {
		t.Errorf("row for foo = %v, want an error message", rows[2])
	}
}

func TestLoggerFillsEmptyFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.csv")
	lg, err := logger.Open(path, "")
	if err != nil {
		t.Fatal(err)
	}
	lg.Log("exit", nil, nil)
	lg.Close()

	row := readCSV(t, path)[1]
	for i, field := range row {
		if field == "" {
			t.Errorf("field %d is empty in row %v, want \"-\"", i, row)
		}
	}
	if row[1] != "-" || row[3] != "-" || row[4] != "-" {
		t.Errorf("row = %v, want \"-\" for user, args and error", row)
	}
}
