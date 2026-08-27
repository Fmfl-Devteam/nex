package network

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var DefaultPorts = []int{22, 53, 80, 443, 3389, 8080}

const maxPorts = 128

type PortResult struct {
	Port  int    `json:"port"`
	Open  bool   `json:"open"`
	Error string `json:"error,omitempty"`
}

// ParsePorts parses a comma-separated, deduplicated list of TCP ports.
func ParsePorts(value string) ([]int, error) {
	if strings.TrimSpace(value) == "" {
		return append([]int(nil), DefaultPorts...), nil
	}
	parts := strings.Split(value, ",")
	if len(parts) > maxPorts {
		return nil, fmt.Errorf("too many ports: maximum is %d", maxPorts)
	}
	seen := make(map[int]struct{}, len(parts))
	ports := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		port, err := strconv.Atoi(part)
		if err != nil || port < 1 || port > 65535 {
			return nil, fmt.Errorf("invalid port %q: must be an integer from 1 to 65535", part)
		}
		if _, exists := seen[port]; !exists {
			seen[port] = struct{}{}
			ports = append(ports, port)
		}
	}
	sort.Ints(ports)
	return ports, nil
}

// CheckPorts attempts bounded concurrent TCP connections and preserves port order.
func CheckPorts(ctx context.Context, host string, ports []int, timeout time.Duration, concurrency int) []PortResult {
	if len(ports) == 0 {
		return []PortResult{}
	}
	if concurrency < 1 {
		concurrency = 1
	}
	if concurrency > len(ports) {
		concurrency = len(ports)
	}
	results := make([]PortResult, len(ports))
	jobs := make(chan int)
	var workers sync.WaitGroup
	for range concurrency {
		workers.Add(1)
		go func() {
			defer workers.Done()
			dialer := net.Dialer{Timeout: timeout}
			for index := range jobs {
				port := ports[index]
				result := PortResult{Port: port}
				conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(host, strconv.Itoa(port)))
				if err != nil {
					result.Error = err.Error()
				} else {
					result.Open = true
					if closeErr := conn.Close(); closeErr != nil {
						result.Error = "close connection: " + closeErr.Error()
					}
				}
				results[index] = result
			}
		}()
	}
	for index := range ports {
		jobs <- index
	}
	close(jobs)
	workers.Wait()
	return results
}
