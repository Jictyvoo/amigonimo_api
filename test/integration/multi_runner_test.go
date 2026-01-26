package integration

import (
	"testing"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners/dbrunner"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners/reqrunner"
)

type UserID string

func TestMultiRunnerStateSharing(t *testing.T) {
	engine := NewEngine(t)

	// Runner 1: Create a user and verify success
	createUserRunner := reqrunner.NewHttpRunner(
		engine.BaseURL(),
		reqrunner.WithRequest(
			"POST", "/auth/sign", controllers.FormUser{
				Username: "stateuser",
				Email:    "state@example.com",
				Password: "password123",
			},
		),
		reqrunner.ExpectStatus(201),
		reqrunner.ExpectBody(
			controllers.SuccessResponse{
				Success: true,
				Message: "User created successfully",
			},
		),
		// Extract email as "UserID" (demo how to pass results)
		reqrunner.ExtractToState(
			func(resp controllers.SuccessResponse) UserID {
				return UserID("state@example.com")
			},
		),
	)

	// Runner 2: Verify in DB that the user exists using the extracted "UserID"
	dbVerifyRunner := dbrunner.NewDbRunner(
		engine.DB(),
		dbrunner.WithSubsequentQuery(
			"users", func(id UserID) map[string]any {
				return map[string]any{"email": string(id)}
			},
		),
		dbrunner.ExpectCount(1),
	)

	mr := runners.MultiRunner{
		Runners: []runners.Runner{
			createUserRunner,
			dbVerifyRunner,
		},
	}

	if err := mr.Run(t); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
