// Package system collects cross-platform host information.
package system

import (
	"fmt"
	"os"
	"runtime"
)

type Info struct {
	OS              string `json:"os"`
	Architecture    string `json:"architecture"`
	Hostname        string `json:"hostname"`
	GoArchitecture  string `json:"go_architecture"`
	CPU             string `json:"cpu,omitempty"`
	LogicalCPUs     int    `json:"logical_cpus"`
	TotalMemory     uint64 `json:"total_memory_bytes,omitempty"`
	AvailableMemory uint64 `json:"available_memory_bytes,omitempty"`
}

// Get collects host details and delegates OS-specific values to build-tagged code.
func Get() (Info, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return Info{}, fmt.Errorf("get hostname: %w", err)
	}
	platform, err := getPlatformInfo()
	if err != nil {
		return Info{}, err
	}
	return Info{
		OS:              platform.OS,
		Architecture:    platform.Architecture,
		Hostname:        hostname,
		GoArchitecture:  runtime.GOARCH,
		CPU:             platform.CPU,
		LogicalCPUs:     runtime.NumCPU(),
		TotalMemory:     platform.TotalMemory,
		AvailableMemory: platform.AvailableMemory,
	}, nil
}

type platformInfo struct {
	OS              string
	Architecture    string
	CPU             string
	TotalMemory     uint64
	AvailableMemory uint64
}
