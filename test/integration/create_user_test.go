package integration

import (
	"net/http"
	"testing"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/bancoche"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	"github.com/jictyvoo/amigonimo_api/test/integration/stdrunners"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
)

func TestCreateUserSimple(t *testing.T) {
	engine := NewEngine(t)

	actorBuilder := fixtures.NewUser().
		WithEmail("actor@example.com")
	actor := actorBuilder.Build()
	actorProfile := actorBuilder.BuildProfile()
	if seedErr := engine.Seed(actor, actorProfile); seedErr != nil {
		t.Fatalf("seedErr: %v", seedErr.Error())
	}

	reqBody := controllers.FormUser{
		Email:    "newuser@example.com",
		Username: "newuser",
		Password: "securepassword",
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodPost, "/auth/sign", reqBody),
				netoche.ExpectStatus(http.StatusCreated),
				netoche.ExpectBody(
					controllers.SuccessResponse{
						Success: true,
						Message: "User created successfully",
					},
				),
			),
			bancoche.New(
				engine.DB(),
				bancoche.WithMapQuery(
					"users", map[string]any{"email": reqBody.Email, "username": reqBody.Username},
				),
				bancoche.ExpectCount(1, true),
			),
			stdrunners.LoginRunner(engine.BaseURL(), reqBody.Email, reqBody.Password),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
