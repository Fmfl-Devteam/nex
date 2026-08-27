// Package proc defines the process command group.
package proc

import (
	"github.com/fmfl-devteam/nex/internal/output"
	"github.com/spf13/cobra"
)

func NewCommand(printer *output.Printer) *cobra.Command {
	command := &cobra.Command{
		Use:   "proc",
		Short: "Inspect and manage processes",
	}
	command.AddCommand(
		newListCommand(printer),
		newInfoCommand(printer),
		newKillCommand(printer),
	)
	return command
}
