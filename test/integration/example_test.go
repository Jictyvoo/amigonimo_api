package integration

import (
	"testing"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners"
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
	
	mr := runners.MultiRunner{
		Runners: []runners.Runner{runner},
	}
	if err := mr.Run(t); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
