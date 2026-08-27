package net

import (
	"fmt"
	"time"

	networking "github.com/fmfl-devteam/nex/internal/network"
	"github.com/fmfl-devteam/nex/internal/output"
	"github.com/spf13/cobra"
)

func newPingCommand(printer *output.Printer) *cobra.Command {
	var timeout time.Duration
	command := &cobra.Command{
		Use:   "ping <host>",
		Short: "Check whether a host responds to ICMP",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			if timeout <= 0 {
				return fmt.Errorf("timeout must be greater than zero")
			}
			result, pingErr := networking.Ping(command.Context(), args[0], timeout)
			var outputErr error
			if printer.JSON {
				outputErr = printer.Encode(result)
			} else {
				status := "reachable"
				if !result.Reachable {
					status = "unreachable"
				}
				outputErr = printer.Table(
					[]string{"HOST", "STATUS", "LATENCY"},
					[]string{result.Host, status, formatLatency(result.LatencyMS)},
				)
			}
			if outputErr != nil {
				return outputErr
			}
			return pingErr
		},
	}
	command.Flags().DurationVar(&timeout, "timeout", 3*time.Second, "maximum time to wait")
	return command
}

func formatLatency(milliseconds float64) string {
	return fmt.Sprintf("%.2f ms", milliseconds)
}
