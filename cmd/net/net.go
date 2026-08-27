// Package net defines the network command group.
package net

import (
	"github.com/fmfl-devteam/nex/internal/output"
	"github.com/spf13/cobra"
)

func NewCommand(printer *output.Printer) *cobra.Command {
	command := &cobra.Command{
		Use:   "net",
		Short: "Run network diagnostics",
	}
	command.AddCommand(
		newPingCommand(printer),
		newDNSCommand(printer),
		newPortsCommand(printer),
	)
	return command
}
