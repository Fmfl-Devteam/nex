package net

import (
	"net"

	networking "github.com/fmfl-devteam/nex/internal/network"
	"github.com/fmfl-devteam/nex/internal/output"
	"github.com/spf13/cobra"
)

func newDNSCommand(printer *output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "dns <host>",
		Short: "Resolve a hostname",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			result, err := networking.Resolve(command.Context(), args[0])
			if err != nil {
				return err
			}
			if printer.JSON {
				return printer.Encode(result)
			}
			rows := [][]string{{"ADDRESS", "TYPE"}}
			for _, address := range result.Addresses {
				family := "IPv6"
				if net.ParseIP(address).To4() != nil {
					family = "IPv4"
				}
				rows = append(rows, []string{address, family})
			}
			return printer.Table(rows...)
		},
	}
}
