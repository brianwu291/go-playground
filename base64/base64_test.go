package base64_test

import (
	"testing"

	"github.com/brianwu291/go-playground/base64"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"empty", []byte(""), ""},
		{"single char", []byte("a"), "YQ=="},
		{"two chars", []byte("ab"), "YWI="},
		{"three chars", []byte("abc"), "YWJj"},
		{"hello world", []byte("Hello, World!"), "SGVsbG8sIFdvcmxkIQ=="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := base64.Encode(tt.input)
			if result != tt.expected {
				t.Errorf("Encode(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []byte
		expectErr bool
	}{
		{"empty", "", []byte{}, false},
		{"single char", "YQ==", []byte("a"), false},
		{"hello world", "SGVsbG8sIFdvcmxkIQ==", []byte("Hello, World!"), false},
		{"invalid length", "ABC", nil, true},
		{"invalid char", "AB@D", nil, true},
		{"invalid char", "好讚11", nil, true},
		{"invalid char", "11好棒", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := base64.Decode(tt.input)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Decode(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Decode(%q) unexpected error: %v", tt.input, err)
				return
			}

			if string(result) != string(tt.expected) {
				t.Errorf("Decode(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	inputs := []string{
		"",
		"a",
		"ab",
		"abc",
		"Hello, World!",
		"The quick brown fox jumps over the lazy dog",
		"1234567890",
		"!@#$%^&*()",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			encoded := base64.Encode([]byte(input))
			decoded, err := base64.Decode(encoded)

			if err != nil {
				t.Errorf("Round trip failed for %q: %v", input, err)
				return
			}

			if string(decoded) != input {
				t.Errorf("Round trip failed: got %q, want %q", decoded, input)
			}
		})
	}
}
