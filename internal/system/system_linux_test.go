//go:build linux

package system

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLinuxMemory(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "meminfo")
	data := []byte("MemTotal:       1024 kB\nMemFree:         100 kB\nMemAvailable:    512 kB\n")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	total, available, err := linuxMemory(path)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1024*1024 || available != 512*1024 {
		t.Fatalf("linuxMemory() = %d, %d", total, available)
	}
}

func TestLinuxCPUModel(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "cpuinfo")
	if err := os.WriteFile(path, []byte("processor: 0\nmodel name: Example CPU\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	model, err := linuxCPUModel(path)
	if err != nil {
		t.Fatal(err)
	}
	if model != "Example CPU" {
		t.Fatalf("linuxCPUModel() = %q", model)
	}
}
