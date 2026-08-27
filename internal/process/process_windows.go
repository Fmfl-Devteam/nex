//go:build windows

package process

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// List uses tasklist, which is included with supported Windows releases.
func List() ([]Summary, error) {
	command := exec.Command("tasklist.exe", "/FO", "CSV", "/NH")
	data, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("run tasklist: %w", err)
	}
	records, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse tasklist output: %w", err)
	}
	items := make([]Summary, 0, len(records))
	for _, record := range records {
		if len(record) < 2 {
			continue
		}
		pid, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			continue
		}
		items = append(items, Summary{PID: pid, Name: record[0]})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].PID < items[j].PID })
	return items, nil
}

// Get uses CIM through Windows PowerShell for fields available to the caller.
func Get(pid int) (Info, error) {
	script := fmt.Sprintf(`$p=Get-CimInstance Win32_Process -Filter "ProcessId=%d"; if ($null -eq $p) { exit 3 }; $o=Invoke-CimMethod -InputObject $p -MethodName GetOwner -ErrorAction SilentlyContinue; [pscustomobject]@{Name=$p.Name;Executable=$p.ExecutablePath;CommandLine=$p.CommandLine;User=if($o.User){if($o.Domain){$o.Domain+'\\'+$o.User}else{$o.User}}else{''}} | ConvertTo-Json -Compress`, pid)
	command := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script)
	data, err := command.Output()
	if err != nil {
		return Info{}, fmt.Errorf("inspect process %d: %w", pid, err)
	}
	var raw struct {
		Name        string
		Executable  string
		CommandLine string
		User        string
	}
	if err := json.Unmarshal(bytes.TrimSpace(data), &raw); err != nil {
		return Info{}, fmt.Errorf("decode process %d information: %w", pid, err)
	}
	result := Info{PID: pid, Name: raw.Name, Executable: raw.Executable, User: raw.User}
	if raw.CommandLine != "" {
		result.Command = []string{raw.CommandLine}
	}
	return result, nil
}
