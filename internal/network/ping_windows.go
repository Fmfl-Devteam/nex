//go:build windows

package network

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

var pingLatencyPattern = regexp.MustCompile(`(?i)time[=<]\s*([0-9.]+)\s*ms`)

// Ping invokes the built-in Windows ping utility with one ICMP request.
func Ping(ctx context.Context, host string, timeout time.Duration) (PingResult, error) {
	milliseconds := timeout.Milliseconds()
	if milliseconds < 1 {
		milliseconds = 1
	}
	started := time.Now()
	data, err := exec.CommandContext(ctx, "ping.exe", "-n", "1", "-w", strconv.FormatInt(milliseconds, 10), host).CombinedOutput()
	elapsed := time.Since(started)
	result := PingResult{Host: host, Reachable: err == nil, Latency: elapsed, LatencyMS: float64(elapsed) / float64(time.Millisecond)}
	if matches := pingLatencyPattern.FindSubmatch(data); len(matches) == 2 {
		if value, parseErr := strconv.ParseFloat(string(matches[1]), 64); parseErr == nil {
			result.LatencyMS = value
			result.Latency = time.Duration(value * float64(time.Millisecond))
		}
	}
	if err != nil {
		return result, fmt.Errorf("ping %q: %w", host, err)
	}
	return result, nil
}
