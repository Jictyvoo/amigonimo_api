package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixturesets"
)

func TestAuthRecoveryRoutes(t *testing.T) {
	engine := NewEngine(t)

	const (
		oldPassword  = "old-password-123"
		newPassword  = "new-password-123"
		recoveryCode = "recover-123"
	)

	user := fixturesets.NewUser("recovery-user@example.com", oldPassword, "").
		WithRecoveryCode(recoveryCode, time.Now().Add(2*time.Hour))

	if err := engine.Seed(user.User, user.Profile); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.CheckRecoveryCode(
				engine.BaseURL(),
				controllers.FormRecoveryCode{
					Email:        user.User.Email,
					RecoveryCode: recoveryCode,
				},
			),
			authrunner.ResetPassword(
				engine.BaseURL(),
				controllers.FormResetPassword{
					Email:        user.User.Email,
					RecoveryCode: recoveryCode,
					NewPassword:  newPassword,
				},
			),
			authrunner.Login(engine.BaseURL(), user.User.Email, newPassword),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestRegenerateRouteRejectsBearerJWT(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "regenerate-password"

	user := fixturesets.NewUser("regenerate-user@example.com", userPassword, "")
	if err := engine.Seed(user.User, user.Profile); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), user.User.Email, userPassword),
			authrunner.RegenerateToken(
				engine.BaseURL(),
				netoche.ExpectStatus(http.StatusPreconditionFailed),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	engine := NewEngine(t)
	const correctPassword = "correct-password-456"

	user := fixturesets.NewUser("wrong-pass-user@example.com", correctPassword, "")
	if err := engine.Seed(user.User, user.Profile); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.FailedLogin(
				engine.BaseURL(),
				user.User.Email,
				"wrong-password",
				http.StatusNotAcceptable,
				"not found user with given email/username and password combination",
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestSignupDuplicateEmail(t *testing.T) {
	engine := NewEngine(t)

	existing := fixturesets.NewUser("duplicate@example.com", "password-xyz", "")
	if err := engine.Seed(existing.User, existing.Profile); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.FailedSignUp(
				engine.BaseURL(),
				controllers.FormUser{
					Email:    existing.User.Email,
					Username: "other-username",
					Password: "another-password",
				},
				http.StatusPreconditionFailed,
				"user already with provided email/username already exists",
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestRecoveryCodeExpired(t *testing.T) {
	engine := NewEngine(t)
	const password = "expired-recovery-pass"

	user := fixturesets.NewUser("expired-recovery@example.com", password, "").
		WithRecoveryCode("expired-code", time.Now().Add(-1*time.Hour)) // already expired

	if err := engine.Seed(user.User, user.Profile); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.FailedCheckRecoveryCode(
				engine.BaseURL(),
				controllers.FormRecoveryCode{
					Email:        user.User.Email,
					RecoveryCode: "expired-code",
				},
				http.StatusPreconditionFailed,
				"cannot find user with given email and recovery code",
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
