package utils

import (
	"strings"
	"testing"
)

func TestKeyValueEntryCompare(t *testing.T) {
	tests := []struct {
		name string
		a    KeyValueEntry[int]
		b    KeyValueEntry[int]
		want int
	}{
		{
			name: "equal keys returns zero",
			a:    KeyValueEntry[int]{Key: "foo"},
			b:    KeyValueEntry[int]{Key: "foo"},
			want: 0,
		},
		{
			name: "a key less than b key returns negative",
			a:    KeyValueEntry[int]{Key: "bar"},
			b:    KeyValueEntry[int]{Key: "foo"},
			want: strings.Compare("bar", "foo"),
		},
		{
			name: "a key greater than b key returns positive",
			a:    KeyValueEntry[int]{Key: "foo"},
			b:    KeyValueEntry[int]{Key: "bar"},
			want: strings.Compare("foo", "bar"),
		},
		{
			name: "empty keys are equal",
			a:    KeyValueEntry[int]{Key: ""},
			b:    KeyValueEntry[int]{Key: ""},
			want: 0,
		},
		{
			name: "value field does not affect comparison",
			a:    KeyValueEntry[int]{Key: "same", Value: 1},
			b:    KeyValueEntry[int]{Key: "same", Value: 99},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Compare(tt.b); got != tt.want {
				t.Errorf("Compare() = %d, want %d", got, tt.want)
			}
		})
	}
}
