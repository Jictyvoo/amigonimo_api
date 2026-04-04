package utils

import (
	"slices"
	"testing"
)

func TestMaxSize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		maxSize int
		want    string
	}{
		{"string shorter than max returned as-is", "hello", 10, "hello"},
		{"string equal to max returned as-is", "hello", 5, "hello"},
		{"string longer than max is cut from start", "hello world", 5, "world"},
		{"zero maxSize defaults to 100", "hello", 0, "hello"},
		{"negative maxSize defaults to 100", "hello", -1, "hello"},
		{"empty string with any max", "", 5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxSize(tt.input, tt.maxSize); got != tt.want {
				t.Errorf("MaxSize(%q, %d) = %q, want %q", tt.input, tt.maxSize, got, tt.want)
			}
		})
	}
}

func TestShuffleString(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"ascii string is permuted", "abcdefgh"},
		{"unicode string is permuted", "αβγδεζηθ"},
		{"single character unchanged", "x"},
		{"empty string unchanged", ""},
	}

	sortRunes := func(s string) string {
		r := []rune(s)
		slices.Sort(r)
		return string(r)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShuffleString(tt.input)
			// Shuffled result must contain the same runes.
			if sortRunes(got) != sortRunes(tt.input) {
				t.Errorf("ShuffleString(%q) = %q, rune set differs", tt.input, got)
			}
			// Length must be preserved.
			if len([]rune(got)) != len([]rune(tt.input)) {
				t.Errorf("ShuffleString(%q) length changed: got %d want %d",
					tt.input, len([]rune(got)), len([]rune(tt.input)))
			}
		})
	}
}
