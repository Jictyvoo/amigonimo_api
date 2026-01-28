package integration

import (
	"net/http"
	"testing"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	"github.com/jictyvoo/amigonimo_api/test/integration/stdrunners"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners/dbrunner"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners/reqrunner"
)

func TestCreateUserSimple(t *testing.T) {
	engine := NewEngine(t)

	actor := fixtures.NewUser().
		WithEmail("actor@example.com").
		Build()
	engine.Seed(actor)

	reqBody := controllers.FormUser{
		Email:    "newuser@example.com",
		Username: "newuser",
		Password: "securepassword",
	}

	mr := runners.MultiRunner{
		Runners: []runners.Runner{
			reqrunner.NewHttpRunner(
				engine.BaseURL(),
				reqrunner.WithRequest(http.MethodPost, "/auth/sign", reqBody),
				reqrunner.ExpectStatus(http.StatusCreated),
				reqrunner.ExpectBody(
					controllers.SuccessResponse{
						Success: true,
						Message: "User created successfully",
					},
				),
			),
			dbrunner.NewDbRunner(
				engine.DB(),
				dbrunner.WithQuery(
					"users", map[string]any{"email": reqBody.Email, "username": reqBody.Username},
				),
				dbrunner.ExpectCount(1),
			),
			stdrunners.LoginRunner(engine.BaseURL(), reqBody.Email, reqBody.Password),
		},
	}
	if err := mr.Run(t); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
