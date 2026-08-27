//go:build linux

package process

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const procRoot = "/proc"

// List returns processes visible through procfs.
func List() ([]Summary, error) {
	entries, err := os.ReadDir(procRoot)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", procRoot, err)
	}
	items := make([]Summary, 0, len(entries))
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || !entry.IsDir() {
			continue
		}
		nameBytes, err := os.ReadFile(filepath.Join(procRoot, entry.Name(), "comm"))
		if err != nil { // Processes can exit while procfs is being read.
			continue
		}
		item := Summary{PID: pid, Name: strings.TrimSpace(string(nameBytes))}
		if info, err := entry.Info(); err == nil {
			item.User = ownerName(info)
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].PID < items[j].PID })
	return items, nil
}

// Get returns details for one process. Permission-protected optional fields are omitted.
func Get(pid int) (Info, error) {
	base := filepath.Join(procRoot, strconv.Itoa(pid))
	nameBytes, err := os.ReadFile(filepath.Join(base, "comm"))
	if err != nil {
		return Info{}, fmt.Errorf("read process %d: %w", pid, err)
	}
	result := Info{PID: pid, Name: strings.TrimSpace(string(nameBytes))}
	if executable, err := os.Readlink(filepath.Join(base, "exe")); err == nil {
		result.Executable = executable
	}
	if cwd, err := os.Readlink(filepath.Join(base, "cwd")); err == nil {
		result.CWD = cwd
	}
	if command, err := os.ReadFile(filepath.Join(base, "cmdline")); err == nil {
		command = bytes.TrimRight(command, "\x00")
		if len(command) > 0 {
			result.Command = strings.Split(string(command), "\x00")
		}
	}
	if stat, err := os.Stat(base); err == nil {
		result.User = ownerName(stat)
	}
	return result, nil
}

func ownerName(info os.FileInfo) string {
	uid := fileUID(info)
	if uid == "" {
		return ""
	}
	account, err := user.LookupId(uid)
	if err != nil {
		return uid
	}
	return account.Username
}
