package network

import "time"

type PingResult struct {
	Host      string        `json:"host"`
	Reachable bool          `json:"reachable"`
	Latency   time.Duration `json:"-"`
	LatencyMS float64       `json:"latency_ms"`
}
