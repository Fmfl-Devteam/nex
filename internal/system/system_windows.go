//go:build windows

package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func getPlatformInfo() (platformInfo, error) {
	script := `$o=Get-CimInstance Win32_OperatingSystem; $c=Get-CimInstance Win32_Processor | Select-Object -First 1; [pscustomobject]@{OS=$o.Caption;CPU=$c.Name;Total=$o.TotalVisibleMemorySize;Free=$o.FreePhysicalMemory} | ConvertTo-Json -Compress`
	data, err := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script).Output()
	if err != nil {
		return platformInfo{}, fmt.Errorf("collect Windows system information: %w", err)
	}
	var raw struct {
		OS    string
		CPU   string
		Total uint64
		Free  uint64
	}
	if err := json.Unmarshal(bytes.TrimSpace(data), &raw); err != nil {
		return platformInfo{}, fmt.Errorf("decode Windows system information: %w", err)
	}
	architecture := strings.TrimSpace(os.Getenv("PROCESSOR_ARCHITECTURE"))
	if architecture == "" {
		architecture = runtime.GOARCH
	}
	return platformInfo{
		OS:              raw.OS,
		Architecture:    architecture,
		CPU:             strings.TrimSpace(raw.CPU),
		TotalMemory:     raw.Total * 1024,
		AvailableMemory: raw.Free * 1024,
	}, nil
}
