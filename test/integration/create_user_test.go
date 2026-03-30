package integration

import (
	"testing"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/bancoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixturesets"
)

func TestCreateUserSimple(t *testing.T) {
	engine := NewEngine(t)

	actor := fixturesets.NewUser("actor@example.com", "password", "")
	if err := engine.Seed(actor.User, actor.Profile); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	reqBody := controllers.FormUser{
		Email:    "newuser@example.com",
		Username: "newuser",
		Password: "securepassword",
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.SignUp(engine.BaseURL(), reqBody),
			bancoche.New(
				engine.DB(),
				bancoche.WithMapQuery(
					"users", map[string]any{"email": reqBody.Email, "username": reqBody.Username},
				),
				bancoche.ExpectCount(1, true),
			),
			bancoche.New(
				engine.DB(),
				bancoche.WithMapQuery("user_profiles", map[string]any{}),
				bancoche.ExpectCount(2, true),
			),
			authrunner.Login(engine.BaseURL(), reqBody.Email, reqBody.Password),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
