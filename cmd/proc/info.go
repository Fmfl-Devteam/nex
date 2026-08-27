package proc

import (
	"strings"

	"github.com/fmfl-devteam/nex/internal/output"
	processes "github.com/fmfl-devteam/nex/internal/process"
	"github.com/spf13/cobra"
)

func newInfoCommand(printer *output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "info <pid>",
		Short: "Show detailed process information",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			pid, err := processes.ParsePID(args[0])
			if err != nil {
				return err
			}
			info, err := processes.Get(pid)
			if err != nil {
				return err
			}
			if printer.JSON {
				return printer.Encode(info)
			}
			return printer.Table(
				[]string{"FIELD", "VALUE"},
				[]string{"PID", args[0]},
				[]string{"Name", info.Name},
				[]string{"Executable", unavailable(info.Executable)},
				[]string{"Working directory", unavailable(info.CWD)},
				[]string{"Command", unavailable(strings.Join(info.Command, " "))},
				[]string{"User", unavailable(info.User)},
			)
		},
	}
}

func unavailable(value string) string {
	if value == "" {
		return "unavailable"
	}
	return value
}
