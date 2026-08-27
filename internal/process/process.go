// Package process exposes cross-platform process inspection operations.
package process

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Summary is the stable, serializable representation used by process lists.
type Summary struct {
	PID  int    `json:"pid"`
	Name string `json:"name"`
	User string `json:"user,omitempty"`
}

// Info contains process details. Empty optional fields mean unavailable.
type Info struct {
	PID        int      `json:"pid"`
	Name       string   `json:"name"`
	Executable string   `json:"executable,omitempty"`
	CWD        string   `json:"cwd,omitempty"`
	Command    []string `json:"command,omitempty"`
	User       string   `json:"user,omitempty"`
}

// ParsePID validates a user-supplied process identifier.
func ParsePID(value string) (int, error) {
	value = strings.TrimSpace(value)
	pid, err := strconv.Atoi(value)
	if err != nil || pid <= 0 {
		return 0, fmt.Errorf("invalid PID %q: must be a positive integer", value)
	}
	return pid, nil
}

// Kill terminates a process using the platform implementation in os.Process.
func Kill(pid int) error {
	if pid <= 0 {
		return fmt.Errorf("invalid PID %d: must be positive", pid)
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("find process %d: %w", pid, err)
	}
	if err := p.Kill(); err != nil {
		return fmt.Errorf("terminate process %d: %w", pid, err)
	}
	return nil
}
