package shell

import (
	"fmt"
	"os"
	"os/user"
)

// unknown replaces a user or host name that cannot be detected.
const unknown = "unknown"

// Title builds the window title from a user name and a host name.
func Title(username, hostname string) string {
	return fmt.Sprintf("Эмулятор - [%s@%s]", username, hostname)
}

// SystemTitle returns the window title for the user and host
// the emulator runs on.
func SystemTitle() string {
	return Title(CurrentUser(), CurrentHost())
}

// CurrentUser returns the name of the OS user running the emulator.
func CurrentUser() string {
	u, err := user.Current()
	if err != nil {
		return unknown
	}
	return u.Username
}

// CurrentHost returns the host name of the machine.
func CurrentHost() string {
	host, err := os.Hostname()
	if err != nil {
		return unknown
	}
	return host
}
