package net

import (
	"fmt"
	"strconv"
	"time"

	networking "github.com/fmfl-devteam/nex/internal/network"
	"github.com/fmfl-devteam/nex/internal/output"
	"github.com/spf13/cobra"
)

func newPortsCommand(printer *output.Printer) *cobra.Command {
	var portsValue string
	var timeout time.Duration
	command := &cobra.Command{
		Use:   "ports <host>",
		Short: "Check a small set of TCP ports",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			if timeout <= 0 {
				return fmt.Errorf("timeout must be greater than zero")
			}
			ports, err := networking.ParsePorts(portsValue)
			if err != nil {
				return err
			}
			results := networking.CheckPorts(command.Context(), args[0], ports, timeout, 32)
			if printer.JSON {
				return printer.Encode(struct {
					Host  string                  `json:"host"`
					Ports []networking.PortResult `json:"ports"`
				}{Host: args[0], Ports: results})
			}
			rows := [][]string{{"PORT", "STATUS", "DETAIL"}}
			for _, result := range results {
				status := "closed/unreachable"
				if result.Open {
					status = "open"
				}
				rows = append(rows, []string{strconv.Itoa(result.Port), status, result.Error})
			}
			return printer.Table(rows...)
		},
	}
	command.Flags().StringVar(&portsValue, "ports", "", "comma-separated TCP ports")
	command.Flags().DurationVar(&timeout, "timeout", 2*time.Second, "connection timeout per port")
	return command
}
