package proc

import (
	"strconv"

	"github.com/fmfl-devteam/nex/internal/output"
	processes "github.com/fmfl-devteam/nex/internal/process"
	"github.com/spf13/cobra"
)

func newListCommand(printer *output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List running processes",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			items, err := processes.List()
			if err != nil {
				return err
			}
			if printer.JSON {
				return printer.Encode(items)
			}
			rows := make([][]string, 0, len(items)+1)
			rows = append(rows, []string{"PID", "NAME", "USER"})
			for _, item := range items {
				rows = append(rows, []string{strconv.Itoa(item.PID), item.Name, item.User})
			}
			return printer.Table(rows...)
		},
	}
}
