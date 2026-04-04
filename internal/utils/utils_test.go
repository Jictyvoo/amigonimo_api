package utils

import "testing"

func TestAbsoluteNum(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"positive value unchanged", 42, 42},
		{"negative value returns positive", -42, 42},
		{"zero unchanged", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbsoluteNum(tt.input); got != tt.want {
				t.Errorf("AbsoluteNum(%d) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestAbsoluteNumFloat(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{"positive float unchanged", 3.14, 3.14},
		{"negative float returns positive", -3.14, 3.14},
		{"zero float unchanged", 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbsoluteNum(tt.input); got != tt.want {
				t.Errorf("AbsoluteNum(%f) = %f, want %f", tt.input, got, tt.want)
			}
		})
	}
}
