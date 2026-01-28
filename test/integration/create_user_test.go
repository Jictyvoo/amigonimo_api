package integration

import (
	"testing"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
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

	runner := reqrunner.NewHttpRunner(
		engine.BaseURL(),
		reqrunner.WithRequest("POST", "/auth/sign", reqBody),
		reqrunner.ExpectStatus(201),
		reqrunner.ExpectBody(
			controllers.SuccessResponse{
				Success: true,
				Message: "User created successfully",
			},
		),
	)

	dbVerifyRunner := dbrunner.NewDbRunner(
		engine.DB(),
		dbrunner.WithQuery(
			"users", map[string]any{"email": reqBody.Email, "username": reqBody.Username},
		),
		dbrunner.ExpectCount(1),
	)

	mr := runners.MultiRunner{
		Runners: []runners.Runner{
			runner,
			dbVerifyRunner,
		},
	}
	if err := mr.Run(t); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
