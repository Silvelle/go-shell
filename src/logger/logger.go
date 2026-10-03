// Package logger writes command events to a CSV file.
package logger

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// filePerm is the permission mode of a newly created log file.
const filePerm = 0o644

// header is the first row of every log file.
var header = []string{"time", "user", "command", "args", "error"}

// Logger appends one CSV row for every executed command.
type Logger struct {
	file   *os.File
	writer *csv.Writer
	user   string
}

// Open opens or creates the CSV log at path. A header row is written
// to a new or empty file. user is recorded in every row.
func Open(path, user string) (*Logger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePerm)
	if err != nil {
		return nil, fmt.Errorf("log: %w", err)
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("log: %w", err)
	}

	l := &Logger{file: file, writer: csv.NewWriter(file), user: user}
	if info.Size() == 0 {
		if err := l.write(header); err != nil {
			file.Close()
			return nil, err
		}
	}
	return l, nil
}

// Log records one command call. cmdErr is the command's error, or nil.
func (l *Logger) Log(command string, args []string, cmdErr error) error {
	message := ""
	if cmdErr != nil {
		message = cmdErr.Error()
	}
	return l.write([]string{
		time.Now().Format(time.RFC3339),
		l.user,
		command,
		strings.Join(args, " "),
		message,
	})
}

// Close closes the log file.
func (l *Logger) Close() error {
	return l.file.Close()
}

// write appends one row and flushes it to disk at once.
func (l *Logger) write(row []string) error {
	if err := l.writer.Write(row); err != nil {
		return fmt.Errorf("log: %w", err)
	}
	l.writer.Flush()
	return l.writer.Error()
}
