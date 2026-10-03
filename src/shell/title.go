package shell

import (
	"fmt"
	"os"
	"os/user"
)

// unknown replaces a user or how  name that cannot be detected.
const unknown = "unknown"

// Title build th window title from a username and a host name.
func Title(username, hostname string) string {
	return fmt.Sprintf("Эмулятор - [%s@%s]", username, hostname)
}

// SystemTitle return the window title for the user and host
// the emulator runs on.
func SystemTitle() string {
	name := unknown
	if u, err := user.Current(); err == nil {
		name = u.Username
	}

	host, err := os.Hostname()
	if err != nil {
		host = unknown
	}
	return Title(name, host)
}
