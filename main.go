package main

import (
	"fmt"
	"os"

	"github.com/fmfl-devteam/nex/cmd"
)

func main() {
	if err := cmd.Execute(os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
