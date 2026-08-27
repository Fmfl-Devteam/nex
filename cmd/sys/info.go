package sys

import (
	"fmt"

	"github.com/fmfl-devteam/nex/internal/output"
	systeminfo "github.com/fmfl-devteam/nex/internal/system"
	"github.com/spf13/cobra"
)

func newInfoCommand(printer *output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show system information",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			info, err := systeminfo.Get()
			if err != nil {
				return err
			}
			if printer.JSON {
				return printer.Encode(info)
			}
			return printer.Table(
				[]string{"FIELD", "VALUE"},
				[]string{"OS", info.OS},
				[]string{"Architecture", info.Architecture},
				[]string{"Hostname", info.Hostname},
				[]string{"Go architecture", info.GoArchitecture},
				[]string{"CPU", info.CPU},
				[]string{"Logical CPUs", fmt.Sprintf("%d", info.LogicalCPUs)},
				[]string{"Total memory", formatBytes(info.TotalMemory)},
				[]string{"Available memory", formatBytes(info.AvailableMemory)},
			)
		},
	}
}

func formatBytes(value uint64) string {
	const unit = uint64(1024)
	if value < unit {
		return fmt.Sprintf("%d B", value)
	}
	divisor, exponent := unit, 0
	for quotient := value / unit; quotient >= unit && exponent < 5; quotient /= unit {
		divisor *= unit
		exponent++
	}
	return fmt.Sprintf("%.1f %ciB", float64(value)/float64(divisor), "KMGTPE"[exponent])
}
