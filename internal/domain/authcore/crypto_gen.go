package authcore

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/hkdf"

	"github.com/jictyvoo/amigonimo_api/internal/utils"
)

const (
	activationTokenSize = 32 // 32 bytes = 256 bits of entropy
	recoveryCodeLength  = 11 // 11 alphanumeric characters
)

// deriveRandomBytes mixes a user-supplied base key with secure random bytes.
// Output size is guaranteed.
func deriveRandomBytes(baseKey string, size int, info []byte) ([]byte, error) {
	// Always start with secure random bytes as the primary entropy
	seed := make([]byte, size)
	if _, err := rand.Read(seed); err != nil {
		// Fallback to UUID if crypto/rand fails
		u, uuidErr := uuid.NewRandom()
		if uuidErr != nil {
			u = uuid.New()
		}
		seed = u[:]
	}

	// HKDF ensures the key is uniformly derived and consistently sized
	h := hkdf.New(sha256.New, seed, []byte(baseKey), info)

	out := make([]byte, size)
	if _, err := io.ReadFull(h, out); err != nil {
		return nil, err
	}

	return out, nil
}

// GenerateActivationToken produces a 100-character URL-safe token.
// The user key and the current timestamp are mixed into HKDF-derived entropy.
// The function is intentionally simple and uses only standard library components.
func GenerateActivationToken(baseKey string) string {
	const outputSize = 100

	// 32 bytes of HKDF-derived high-entropy material
	core, err := deriveRandomBytes(baseKey, activationTokenSize, []byte("activation-token"))
	if err != nil {
		// extremely rare fallback: use time as seed
		core = []byte(time.Now().String())
	}

	// Add a small, unique timestamp salt so tokens never repeat
	timeSalt := []byte(time.Now().UTC().Format(time.RFC3339Nano))

	// Derive a second buffer to enrich entropy (64 bytes is enough for b64)
	var extra []byte
	if extra, err = deriveRandomBytes(
		baseKey, 64, append([]byte("activation-extra"), timeSalt...),
	); err != nil {
		shaSalt := sha256.Sum256(timeSalt)
		extra = shaSalt[:]
	}

	// Combine the two buffers into one large token material
	merged := append(core, extra...)

	// URL-safe base64, no padding
	tok := base64.RawURLEncoding.EncodeToString(merged)
	if len(tok) < outputSize {
		// pad deterministically with the token itself
		return tok + tok[:outputSize-len(tok)]
	}

	return tok[:outputSize]
}

// GenerateRecoveryCode creates an 11-character uppercase alphanumeric code.
// The entropy source is HKDF derived from secure randomness and the base key.
//
// Output uses Base32 without padding, trimmed to a fixed length.
func GenerateRecoveryCode(baseKey string) (string, error) {
	// 16 bytes gives enough Base32 characters to extract 11 usable chars
	const rawSize = 16

	bytes, err := deriveRandomBytes(baseKey, rawSize, []byte("recovery-code"))
	if err != nil {
		return "", fmt.Errorf("failed to derive random bytes: %w", err)
	}

	// Base32 is uppercase A-Z + 2-7, safe for manual entry
	code := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)

	// Enforce length exactly
	if len(code) < recoveryCodeLength {
		// Should not happen, but enforce defensive pad
		code = utils.ShuffleString(code + code)
	}

	return code[:recoveryCodeLength], nil
}
