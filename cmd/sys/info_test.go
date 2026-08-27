package sys

import "testing"

func TestFormatBytes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		value uint64
		want  string
	}{
		{value: 0, want: "0 B"},
		{value: 1023, want: "1023 B"},
		{value: 1024, want: "1.0 KiB"},
		{value: 1536, want: "1.5 KiB"},
		{value: 1024 * 1024, want: "1.0 MiB"},
	}
	for _, test := range tests {
		if got := formatBytes(test.value); got != test.want {
			t.Errorf("formatBytes(%d) = %q, want %q", test.value, got, test.want)
		}
	}
}
