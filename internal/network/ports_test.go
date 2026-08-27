package network

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestParsePorts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		value   string
		want    []int
		wantErr bool
	}{
		{name: "defaults", value: "", want: DefaultPorts},
		{name: "sort and deduplicate", value: "443, 80,443,22", want: []int{22, 80, 443}},
		{name: "minimum", value: "1", want: []int{1}},
		{name: "maximum", value: "65535", want: []int{65535}},
		{name: "zero", value: "0", wantErr: true},
		{name: "too large", value: "65536", wantErr: true},
		{name: "missing entry", value: "80,,443", wantErr: true},
		{name: "not numeric", value: "http", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParsePorts(test.value)
			if (err != nil) != test.wantErr {
				t.Fatalf("ParsePorts(%q) error = %v, wantErr %v", test.value, err, test.wantErr)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ParsePorts(%q) = %v, want %v", test.value, got, test.want)
			}
		})
	}
}

func TestCheckPortsDetectsListener(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	results := CheckPorts(context.Background(), "127.0.0.1", []int{port}, time.Second, 4)
	if len(results) != 1 || !results[0].Open || results[0].Port != port {
		t.Fatalf("CheckPorts() = %+v, want open port %d", results, port)
	}
}

func TestCheckPortsEmptyInput(t *testing.T) {
	t.Parallel()
	results := CheckPorts(context.Background(), "localhost", nil, time.Second, 4)
	if results == nil || len(results) != 0 {
		t.Fatalf("CheckPorts() = %#v, want non-nil empty result", results)
	}
}
