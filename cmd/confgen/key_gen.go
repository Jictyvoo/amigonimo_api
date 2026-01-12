package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func generateAuthKey() string {
	const rsaKeySize = 2048
	// Generate a 2048-bit RSA private key for PS256 JWT signing
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		// Fallback: if generation fails, return empty string
		// The config loader should handle this case
		return ""
	}

	// Encode the private key in PKCS8 format (preferred format)
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		// Fallback to PKCS1 format if PKCS8 fails
		privateKeyBytes = x509.MarshalPKCS1PrivateKey(privateKey)
	}

	// Create PEM block
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)

	return string(privateKeyPEM)
}
