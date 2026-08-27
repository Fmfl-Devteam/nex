//go:build linux

package system

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

func getPlatformInfo() (platformInfo, error) {
	result := platformInfo{OS: "linux", Architecture: linuxArchitecture()}
	cpu, err := linuxCPUModel("/proc/cpuinfo")
	if err != nil {
		return platformInfo{}, err
	}
	result.CPU = cpu
	total, available, err := linuxMemory("/proc/meminfo")
	if err != nil {
		return platformInfo{}, err
	}
	result.TotalMemory = total
	result.AvailableMemory = available
	return result, nil
}

func linuxArchitecture() string {
	var name syscall.Utsname
	if err := syscall.Uname(&name); err != nil {
		return runtime.GOARCH
	}
	bytes := make([]byte, 0, len(name.Machine))
	for _, value := range name.Machine {
		if value == 0 {
			break
		}
		bytes = append(bytes, byte(value))
	}
	return string(bytes)
}

func linuxCPUModel(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("read CPU information: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, found := strings.Cut(scanner.Text(), ":")
		if found && (strings.TrimSpace(key) == "model name" || strings.TrimSpace(key) == "Hardware") {
			return strings.TrimSpace(value), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scan CPU information: %w", err)
	}
	return "", nil
}

func linuxMemory(path string) (uint64, uint64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, fmt.Errorf("read memory information: %w", err)
	}
	defer file.Close()
	values := make(map[string]uint64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		if key != "MemTotal" && key != "MemAvailable" && key != "MemFree" {
			continue
		}
		kilobytes, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("parse %s: %w", key, err)
		}
		values[key] = kilobytes * 1024
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, fmt.Errorf("scan memory information: %w", err)
	}
	available := values["MemAvailable"]
	if available == 0 {
		available = values["MemFree"]
	}
	return values["MemTotal"], available, nil
}
