package authvalues

import (
	"strings"
	"testing"
)

func TestUserBasicObfuscateEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		email string
		want  string
	}{
		{
			name:  "long local part keeps prefix and masks alternating tail",
			email: "johndoe@example.com",
			want:  "joh**@example.com",
		},
		{
			name:  "short local part keeps only first character",
			email: "ab@example.com",
			want:  "a@example.com",
		},
		{
			name:  "missing host keeps separator behavior stable",
			email: "user",
			want:  "us*@",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user := UserBasic{Email: tt.email}

			if got := user.ObfuscateEmail(); got != tt.want {
				t.Fatalf("ObfuscateEmail() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUserBasicEncryptPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "hashes a regular password",
			password: "super-secret",
		},
		{
			name:     "hashes an empty password",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user := UserBasic{Password: tt.password}

			got, err := user.EncryptPassword()
			if err != nil {
				t.Fatalf("EncryptPassword() error = %v, want nil", err)
			}
			if len(got) == 0 {
				t.Fatal("EncryptPassword() returned empty hash")
			}
			if tt.password != "" && strings.Contains(string(got), tt.password) {
				t.Fatal("EncryptPassword() leaked plaintext password")
			}
		})
	}
}

func TestUserBasicComparePassword(t *testing.T) {
	t.Parallel()

	hashedPassword, err := UserBasic{Password: "super-secret"}.EncryptPassword()
	if err != nil {
		t.Fatalf("EncryptPassword() error = %v, want nil", err)
	}

	tests := []struct {
		name      string
		password  string
		otherPass string
		wantOK    bool
		wantErr   bool
	}{
		{
			name:      "returns true for matching password",
			password:  string(hashedPassword),
			otherPass: "super-secret",
			wantOK:    true,
		},
		{
			name:      "returns false for mismatched password",
			password:  string(hashedPassword),
			otherPass: "wrong-password",
			wantOK:    false,
			wantErr:   true,
		},
		{
			name:      "returns false for invalid hash",
			password:  "not-a-bcrypt-hash",
			otherPass: "super-secret",
			wantOK:    false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user := UserBasic{Password: tt.password}

			gotOK, err := user.ComparePassword(tt.otherPass)
			if gotOK != tt.wantOK {
				t.Fatalf("ComparePassword() ok = %v, want %v", gotOK, tt.wantOK)
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("ComparePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
