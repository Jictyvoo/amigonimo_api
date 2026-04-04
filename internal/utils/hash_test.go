package utils

import (
	"testing"
)

func TestHash(t *testing.T) {
	type simple struct{ Name string }

	tests := []struct {
		name        string
		a           any
		b           any
		wantSameErr bool // both calls should error
	}{
		{
			name: "equal structs produce equal hashes",
			a:    simple{Name: "hello"},
			b:    simple{Name: "hello"},
		},
		{
			name: "different structs produce different hashes",
			a:    simple{Name: "hello"},
			b:    simple{Name: "world"},
		},
		{
			name: "integer values hash consistently",
			a:    42,
			b:    42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashA, errA := Hash(tt.a)
			hashB, errB := Hash(tt.b)

			if errA != nil {
				t.Fatalf("Hash(a) unexpected error: %v", errA)
			}
			if errB != nil {
				t.Fatalf("Hash(b) unexpected error: %v", errB)
			}

			if tt.a == tt.b {
				if hashA != hashB {
					t.Errorf("equal inputs: Hash(a)=%q != Hash(b)=%q", hashA, hashB)
				}
			} else {
				if hashA == hashB {
					t.Errorf("different inputs: Hash(a)=%q == Hash(b)=%q (collision)", hashA, hashB)
				}
			}
		})
	}
}

func TestHashIdempotent(t *testing.T) {
	type payload struct{ ID int }
	v := payload{ID: 7}

	first, err := Hash(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := Hash(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first != second {
		t.Errorf("Hash is not idempotent: %q != %q", first, second)
	}
}
