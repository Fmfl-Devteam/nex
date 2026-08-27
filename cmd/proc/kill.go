package proc

import (
	"fmt"
	"strconv"

	"github.com/fmfl-devteam/nex/internal/output"
	processes "github.com/fmfl-devteam/nex/internal/process"
	"github.com/spf13/cobra"
)

func newKillCommand(printer *output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "kill <pid>",
		Short: "Terminate a process",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			pid, err := processes.ParsePID(args[0])
			if err != nil {
				return err
			}
			if err := processes.Kill(pid); err != nil {
				return err
			}
			if printer.JSON {
				return printer.Encode(struct {
					PID    int    `json:"pid"`
					Status string `json:"status"`
				}{PID: pid, Status: "terminated"})
			}
			_, err = fmt.Fprintf(printer.Writer, "Terminated process %s\n", strconv.Itoa(pid))
			if err != nil {
				return fmt.Errorf("write result: %w", err)
			}
			return nil
		},
	}
}
