package protocol

import (
	"bufio"
	"strings"
	"testing"
)

func TestSerialize(t *testing.T) {
	tests := []struct {
		name     string
		input    RESPValue
		expected string
	}{
		{
			name:     "simple string",
			input:    RESPValue{Type: SimpleString, Str: "OK"},
			expected: "+OK\r\n",
		},
		{
			name:     "error",
			input:    RESPValue{Type: Error, Str: "Error"},
			expected: "-Error message\r\n",
		},
		{
			name:     "integer",
			input:    RESPValue{Type: Integer, Num: 123},
			expected: ":123\r\n",
		},
		{
			name:     "bulk string",
			input:    RESPValue{Type: BulkString, Str: "hello ak"},
			expected: "$8\r\nhello ak\r\n",
		},
		{
			name: "array",
			input: RESPValue{Type: Array, Array: []RESPValue{
				{Type: BulkString, Str: "testing"},
				{Type: BulkString, Str: "the"},
				{Type: BulkString, Str: "array"},
			}},
			expected: "*3\r\n$7\r\ntesting\r\n$3\r\nthe\r\n$5\r\narray\r\n",
		},
		{
			name:     "null",
			input:    RESPValue{Type: BulkString, IsNull: true},
			expected: "$-1\r\n",
		},
		{
			name: "array_2",
			input: RESPValue{Type: Array, Array: []RESPValue{
				{Type: BulkString, Str: "ping"},
			}},
			expected: "*1\r\n$4\r\nping\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Serialize()
			if result != tt.expected {
				t.Errorf("Got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDeserialize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected RESPValue
	}{
		{
			name:     "simple_string",
			input:    "+OK\r\n",
			expected: RESPValue{Type: SimpleString, Str: "OK"},
		},
		{
			name:     "error",
			input:    "-Error message\r\n",
			expected: RESPValue{Type: Error, Str: "Error message"},
		},
		{
			name:     "integer",
			input:    ":123\r\n",
			expected: RESPValue{Type: Integer, Num: 123},
		},
		{
			name:     "bulk_string",
			input:    "$8\r\nhello ak\r\n",
			expected: RESPValue{Type: BulkString, Str: "hello ak"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			result, err := Deserialize(reader)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("Received type %c, want type %c", result.Type, tt.expected.Type)
			}
			if result.Str != tt.expected.Str {
				t.Errorf("Received string %q, want string %q", result.Str, tt.expected.Str)
			}
		})
	}
}
