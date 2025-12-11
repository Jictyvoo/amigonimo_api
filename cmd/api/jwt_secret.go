package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/wrapped-owls/goremy-di/remy"
)

func registerSecret(secretKey []byte, inj remy.Injector) error {
	// Parse RSA private key from PEM-encoded bytes
	block, _ := pem.Decode(secretKey)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block from secret key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1 format as fallback
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("key is not an RSA private key")
	}
	// Register the parsed RSA key for injection
	remy.RegisterInstance(inj, rsaKey)
	return nil
}
