package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestEncodeProducesValidJSON(t *testing.T) {
	t.Parallel()
	var buffer bytes.Buffer
	printer := New(&buffer)
	value := struct {
		Name string `json:"name"`
	}{Name: "nex"}
	if err := printer.Encode(value); err != nil {
		t.Fatal(err)
	}
	var decoded map[string]string
	if err := json.Unmarshal(buffer.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON %q: %v", buffer.String(), err)
	}
	if decoded["name"] != "nex" {
		t.Fatalf("decoded name = %q", decoded["name"])
	}
}

func TestTableAlignsColumns(t *testing.T) {
	t.Parallel()
	var buffer bytes.Buffer
	if err := New(&buffer).Table([]string{"PID", "NAME"}, []string{"1", "init"}); err != nil {
		t.Fatal(err)
	}
	if got := buffer.String(); !strings.Contains(got, "PID  NAME") || !strings.Contains(got, "1    init") {
		t.Fatalf("unexpected table output %q", got)
	}
}
