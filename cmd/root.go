package cmd

import (
	"io"

	netcmd "github.com/fmfl-devteam/nex/cmd/net"
	proccmd "github.com/fmfl-devteam/nex/cmd/proc"
	syscmd "github.com/fmfl-devteam/nex/cmd/sys"
	"github.com/fmfl-devteam/nex/internal/output"
	"github.com/spf13/cobra"
)

// NewCommand constructs the complete nex command tree.
func NewCommand(stdout, stderr io.Writer) *cobra.Command {
	printer := output.New(stdout)
	root := &cobra.Command{
		Use:           "nex",
		Short:         "Cross-platform system and network utilities",
		Long:          "Nex is a cross-platform collection of process, network, and system diagnostic utilities.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.PersistentFlags().BoolVar(&printer.JSON, "json", false, "output results as JSON")
	root.AddCommand(
		proccmd.NewCommand(printer),
		netcmd.NewCommand(printer),
		syscmd.NewCommand(printer),
	)
	return root
}

// Execute runs nex with process-standard streams.
func Execute(stdout, stderr io.Writer) error {
	return NewCommand(stdout, stderr).Execute()
}
