package shell

import (
	"reflect"
	"testing"
)

type test struct {
	name     string
	input    string
	wantName string
	wantArgs string
}

func TestParse(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		wantName string
		wantArgs []string
	}{
		{
			name:     "empty line",
			input:    "",
			wantName: "",
			wantArgs: []string{},
		},
		{
			name:     "whitespace only",
			input:    "    \t   ",
			wantName: "",
			wantArgs: []string{},
		},
		{
			name:     "command without arguments",
			input:    "ls",
			wantName: "ls",
			wantArgs: []string{},
		},
		{
			name:     "command with arguments",
			input:    "ls  -l /home",
			wantName: "ls",
			wantArgs: []string{"-l", "/home"},
		},
		{
			name:     "vsurroding and repeated spaces are collapsed",
			input:    "cd /tmp",
			wantName: "cd",
			wantArgs: []string{"/tmp"},
		},
		{
			name:     "tabs separete tokens like spaces",
			input:    "ls\t-a",
			wantName: "ls",
			wantArgs: []string{"-a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.Name != tt.wantName {
				t.Errorf("Parse().Name = %v, want %v", got.Name, tt.wantName)
			}

			if got.Args == nil {
				t.Fatalf("Parse(%q).Args is nil, want a non-nil slice", tt.input)
			}
			if !reflect.DeepEqual(got.Args, tt.wantArgs) {
				t.Errorf("Parse(%q).Args = %#v, want %#v", tt.input, got.Args, tt.wantArgs)
			}
		})
	}
}

func TestCommandIsEmpty(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{input: "ls", want: false},
		{input: "", want: true},
		{input: "       ", want: true},
		{input: " cd /tmp", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.IsEmpty() != tt.want {
				t.Errorf("Parse().IsEmpty() = %v, want %v", got.IsEmpty(), tt.want)
			}
		})
	}
}
