package testengine

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http/httptest"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/bootstrap"
)

func WithTestServer(t testing.TB, e *Engine) {
	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	remy.RegisterInstance(inj, e.db)

	conf := bootstrap.Config()
	bootstrap.DoInjections(inj, conf)

	const rsaKeySize = 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		t.Fatalf("failed to generate private key: %s", err.Error())
	}

	remy.RegisterInstance(inj, privateKey)
	server, servErr := bootstrap.NewWebServer(conf, &privateKey.PublicKey, inj)
	if servErr != nil {
		t.Fatalf("failed to start server: %s", servErr.Error())
	}

	t.Cleanup(
		func() {
			if closeErr := server.Close(); err != nil {
				t.Errorf("failed to close server: %s", closeErr.Error())
			}
		},
	)

	// Check if server.Mux exists or use server directly if it satisfies Handler.
	// Previously we assumed server.Mux.
	e.ts = httptest.NewServer(server.Mux)
}
