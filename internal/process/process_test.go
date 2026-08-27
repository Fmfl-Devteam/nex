package process

import "testing"

func TestParsePID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		value   string
		want    int
		wantErr bool
	}{
		{name: "valid", value: "1234", want: 1234},
		{name: "surrounding whitespace", value: " 42 ", want: 42},
		{name: "zero", value: "0", wantErr: true},
		{name: "negative", value: "-9", wantErr: true},
		{name: "text", value: "abc", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pid, err := ParsePID(test.value)
			if (err != nil) != test.wantErr {
				t.Fatalf("ParsePID(%q) error = %v, wantErr %v", test.value, err, test.wantErr)
			}
			if pid != test.want {
				t.Errorf("ParsePID(%q) = %d, want %d", test.value, pid, test.want)
			}
		})
	}
}
