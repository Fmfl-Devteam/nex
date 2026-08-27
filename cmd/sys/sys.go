// Package sys defines the system command group.
package sys

import (
	"github.com/fmfl-devteam/nex/internal/output"
	"github.com/spf13/cobra"
)

func NewCommand(printer *output.Printer) *cobra.Command {
	command := &cobra.Command{
		Use:   "sys",
		Short: "Inspect the local system",
	}
	command.AddCommand(newInfoCommand(printer))
	return command
}
