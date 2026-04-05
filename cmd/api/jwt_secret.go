package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/wrapped-owls/goremy-di/remy"
)

func registerSecret(secretKey []byte, inj remy.Injector) (*rsa.PublicKey, error) {
	// Docker (and some env var sources) may carry literal `\n` (two chars: backslash+n)
	// instead of real newline bytes. Normalise both variants so pem.Decode can find the block.
	normalised := bytes.ReplaceAll(secretKey, []byte(`\n`), []byte("\n"))

	// Parse RSA private key from PEM-encoded bytes
	block, _ := pem.Decode(normalised)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from secret key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1 format as fallback
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA private key")
	}
	// Register the parsed RSA private key for injection
	remy.RegisterInstance(inj, rsaKey)
	// Extract and return the public key for JWT verification
	publicKey := &rsaKey.PublicKey
	return publicKey, nil
}
