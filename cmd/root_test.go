package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootShowsHelpWithoutOperation(t *testing.T) {
	var stdout, stderr bytes.Buffer
	command := NewCommand(&stdout, &stderr)
	command.SetArgs(nil)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	output := stdout.String()
	for _, expected := range []string{"Usage:", "Available Commands:", "proc", "net", "sys", "Flags:", "help"} {
		if !strings.Contains(output, expected) {
			t.Errorf("help output does not contain %q:\n%s", expected, output)
		}
	}
}
