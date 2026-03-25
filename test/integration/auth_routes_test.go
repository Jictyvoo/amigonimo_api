package integration

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	"github.com/jictyvoo/amigonimo_api/test/integration/stdrunners"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
)

func TestAuthRecoveryRoutes(t *testing.T) {
	engine := NewEngine(t)

	const (
		oldPassword  = "old-password-123"
		newPassword  = "new-password-123"
		recoveryCode = "recover-123"
	)

	userBuilder := fixtures.NewUser().
		WithEmail("recovery-user@example.com").
		WithPassword(oldPassword)
	user := userBuilder.Build()
	user.RecoveryCode = sql.NullString{String: recoveryCode, Valid: true}
	user.RecoveryCodeExpiresAt = sql.NullTime{
		Time:  time.Now().Add(2 * time.Hour),
		Valid: true,
	}
	userProfile := userBuilder.BuildProfile()

	if err := engine.Seed(user, userProfile); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(
					http.MethodPatch,
					"/auth/password/check-recovery",
					controllers.FormRecoveryCode{
						Email:        user.Email,
						RecoveryCode: recoveryCode,
					},
				),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					controllers.SuccessResponse{
						Success: true,
						Message: "sent recovery-code is valid",
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(
					http.MethodPut,
					"/auth/password/reset",
					controllers.FormResetPassword{
						Email:        user.Email,
						RecoveryCode: recoveryCode,
						NewPassword:  newPassword,
					},
				),
				netoche.ExpectStatus(http.StatusCreated),
				netoche.ExpectBody(
					controllers.SuccessResponse{
						Success: true,
						Message: "password changed",
					},
				),
			),
			stdrunners.LoginRunner(engine.BaseURL(), user.Email, newPassword),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestRegenerateRouteRejectsBearerJWT(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "regenerate-password"

	userBuilder := fixtures.NewUser().
		WithEmail("regenerate-user@example.com").
		WithPassword(userPassword)
	user := userBuilder.Build()
	userProfile := userBuilder.BuildProfile()

	if err := engine.Seed(user, userProfile); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			stdrunners.LoginRunner(engine.BaseURL(), user.Email, userPassword),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodPatch, "/auth/regenerate", struct{}{}),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.ExpectStatus(http.StatusPreconditionFailed),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
