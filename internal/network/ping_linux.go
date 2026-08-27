//go:build linux

package network

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

var pingLatencyPattern = regexp.MustCompile(`time[=<]([0-9.]+)\s*ms`)

// Ping invokes the unprivileged system ping utility with one ICMP request.
func Ping(ctx context.Context, host string, timeout time.Duration) (PingResult, error) {
	seconds := int((timeout + time.Second - 1) / time.Second)
	started := time.Now()
	data, err := exec.CommandContext(ctx, "ping", "-c", "1", "-W", strconv.Itoa(seconds), host).CombinedOutput()
	elapsed := time.Since(started)
	result := PingResult{Host: host, Reachable: err == nil, Latency: elapsed, LatencyMS: float64(elapsed) / float64(time.Millisecond)}
	if matches := pingLatencyPattern.FindSubmatch(data); len(matches) == 2 {
		if milliseconds, parseErr := strconv.ParseFloat(string(matches[1]), 64); parseErr == nil {
			result.LatencyMS = milliseconds
			result.Latency = time.Duration(milliseconds * float64(time.Millisecond))
		}
	}
	if err != nil {
		return result, fmt.Errorf("ping %q: %w", host, err)
	}
	return result, nil
}
