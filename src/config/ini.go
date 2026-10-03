package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// ErrINISyntax is returned for a line that is not "key = value".
var ErrINISyntax = errors.New("invalid INI line")

// ErrUnknownKey ir returned for a key the emulator does not support.
var ErrUnknownKey = errors.New("unknown config key")

// ReadINI reads key = value pairs from the INI file at path.
func ReadINI(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return ParseINI(string(data))
}

// ParseINI parses INI text into key = value pairs.
// Blank lines, comments (; or #) and [section] headers are ignored.
func ParseINI(text string) (map[string]string, error) {
	values := map[string]string{}
	for i, raw := range strings.Split(text, "\n") {
		line := strings.TrimSpace(raw)
		if isSkippable(line) {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("%w: line %d: %q", ErrINISyntax, i+1, line)
		}
		values[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return values, nil
}

// isSkippable reports whether an INI line carries no key = value pair.
func isSkippable(line string) bool {
	isComment := strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#")
	isSection := strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")
	return line == "" || isComment || isSection
}
