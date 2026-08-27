// Package network contains platform-independent network diagnostics.
package network

import (
	"context"
	"fmt"
	"net"
	"sort"
)

type DNSResult struct {
	Host      string   `json:"host"`
	Addresses []string `json:"addresses"`
}

// Resolve looks up and deterministically sorts all IP addresses for host.
func Resolve(ctx context.Context, host string) (DNSResult, error) {
	addresses, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return DNSResult{}, fmt.Errorf("resolve %q: %w", host, err)
	}
	seen := make(map[string]struct{}, len(addresses))
	result := DNSResult{Host: host, Addresses: make([]string, 0, len(addresses))}
	for _, address := range addresses {
		value := address.IP.String()
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result.Addresses = append(result.Addresses, value)
	}
	sort.Slice(result.Addresses, func(i, j int) bool {
		left, right := net.ParseIP(result.Addresses[i]), net.ParseIP(result.Addresses[j])
		leftV4, rightV4 := left.To4() != nil, right.To4() != nil
		if leftV4 != rightV4 {
			return leftV4
		}
		return result.Addresses[i] < result.Addresses[j]
	})
	return result, nil
}
