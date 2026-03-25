package autherrs

import (
	"errors"
	"net/http"
	"testing"
)

func assertErrorShape(
	t *testing.T,
	err *Error,
	wantCode string,
	wantMsg string,
	wantStatus int,
) {
	t.Helper()

	if err.Error() != wantMsg {
		t.Fatalf("Error() = %q, want %q", err.Error(), wantMsg)
	}
	if err.DetailMsg() != wantMsg {
		t.Fatalf("DetailMsg() = %q, want %q", err.DetailMsg(), wantMsg)
	}
	if err.Code() != wantCode {
		t.Fatalf("Code() = %q, want %q", err.Code(), wantCode)
	}
	if err.StatusCode() != wantStatus {
		t.Fatalf("StatusCode() = %d, want %d", err.StatusCode(), wantStatus)
	}
}

func assertErrorExtractors(t *testing.T, err *Error) {
	t.Helper()

	got, ok := errors.AsType[*Error](err)
	if !ok {
		t.Fatal("errors.AsType[*Error](err) = false, want true")
	}
	if !errors.Is(got, err) {
		t.Fatal("errors.AsType[*Error](err) did not return an equivalent error")
	}

	var gotErr *Error
	if !errors.As(err, &gotErr) {
		t.Fatal("errors.As(err, &gotErr) = false, want true")
	}
	if !errors.Is(gotErr, err) {
		t.Fatal("errors.As(err, &gotErr) did not return an equivalent error")
	}
}

func assertErrorIs(t *testing.T, err *Error, wantCode string, wantStatus int) {
	t.Helper()

	if !errors.Is(err, err) {
		t.Fatal("errors.Is(err, err) = false, want true")
	}
	if !errors.Is(err, newError(wantCode, wantStatus, "other", nil, nil)) {
		t.Fatal("errors.Is(err, same-code target) = false, want true")
	}
}

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        *Error
		wantCode   string
		wantMsg    string
		wantStatus int
	}{
		{
			name:       "email or username used",
			err:        ErrEmailOrUsernameUsed,
			wantCode:   "auth_email_or_username_used",
			wantMsg:    "user already with provided email/username already exists",
			wantStatus: http.StatusPreconditionFailed,
		},
		{
			name:       "invalid credentials",
			err:        ErrInvalidCredentials,
			wantCode:   "auth_invalid_credentials",
			wantMsg:    "not found user with given email/username and password combination",
			wantStatus: http.StatusNotAcceptable,
		},
		{
			name:       "wrong password",
			err:        ErrWrongPassword,
			wantCode:   "auth_wrong_password",
			wantMsg:    "provided password don't match with password for this user",
			wantStatus: http.StatusNotAcceptable,
		},
		{
			name:       "verification code",
			err:        ErrVerificationCode,
			wantCode:   "auth_verification_code_invalid",
			wantMsg:    "cannot find any user with given code",
			wantStatus: http.StatusPreconditionFailed,
		},
		{
			name:       "user email not found",
			err:        ErrUserEmailNotFound,
			wantCode:   "auth_user_email_not_found",
			wantMsg:    "user not found",
			wantStatus: http.StatusNotAcceptable,
		},
		{
			name:       "user recovery not found",
			err:        ErrUserRecoveryNotFound,
			wantCode:   "auth_user_recovery_not_found",
			wantMsg:    "cannot find user with given email and recovery code",
			wantStatus: http.StatusPreconditionFailed,
		},
		{
			name:       "user not found",
			err:        ErrUserNotFound,
			wantCode:   "auth_user_not_found",
			wantMsg:    "user not found, make sure you have a valid session",
			wantStatus: http.StatusNotAcceptable,
		},
		{
			name:       "invalid auth token",
			err:        ErrInvalidAuthToken,
			wantCode:   "auth_invalid_token",
			wantMsg:    "error token provided is invalid",
			wantStatus: http.StatusPreconditionFailed,
		},
		{
			name:       "email in use",
			err:        ErrEmailInUse,
			wantCode:   "auth_email_in_use",
			wantMsg:    "this email is already in use and cannot be changed",
			wantStatus: http.StatusNotAcceptable,
		},
		{
			name:       "username in use",
			err:        ErrUsernameInUse,
			wantCode:   "auth_username_in_use",
			wantMsg:    "this username is already in use and cannot be changed",
			wantStatus: http.StatusNotAcceptable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertErrorShape(t, tt.err, tt.wantCode, tt.wantMsg, tt.wantStatus)
			assertErrorExtractors(t, tt.err)
			assertErrorIs(t, tt.err, tt.wantCode, tt.wantStatus)
		})
	}
}

func TestWrappedConstructors(t *testing.T) {
	internalErr := errors.New("sql: duplicate key value violates constraint")

	tests := []struct {
		name       string
		err        *Error
		wantCode   string
		wantMsg    string
		wantStatus int
		wantUnwrap error
	}{
		{
			name:       "password encryption",
			err:        NewErrPasswordEncryption(internalErr),
			wantCode:   "auth_password_encryption_failed",
			wantMsg:    "password encryption error",
			wantStatus: http.StatusInternalServerError,
			wantUnwrap: internalErr,
		},
		{
			name:       "update password",
			err:        NewErrUpdatePassword(internalErr),
			wantCode:   "auth_update_password_failed",
			wantMsg:    "internal error: cannot update the password",
			wantStatus: http.StatusInternalServerError,
			wantUnwrap: internalErr,
		},
		{
			name:       "update username",
			err:        NewErrUpdateUsername(internalErr),
			wantCode:   "auth_update_username_failed",
			wantMsg:    "internal error: cannot update the username",
			wantStatus: http.StatusInternalServerError,
			wantUnwrap: internalErr,
		},
		{
			name:       "update email",
			err:        NewErrUpdateEmail(internalErr),
			wantCode:   "auth_update_email_failed",
			wantMsg:    "internal error: cannot update the email",
			wantStatus: http.StatusInternalServerError,
			wantUnwrap: internalErr,
		},
		{
			name:       "user creation",
			err:        NewErrUserCreation(internalErr),
			wantCode:   "auth_user_creation_failed",
			wantMsg:    "user can't be created",
			wantStatus: http.StatusInternalServerError,
			wantUnwrap: internalErr,
		},
		{
			name:       "update auth token",
			err:        NewErrUpdateAuthToken(internalErr),
			wantCode:   "auth_update_token_failed",
			wantMsg:    "was not possible to update the user authentication token",
			wantStatus: http.StatusInternalServerError,
			wantUnwrap: internalErr,
		},
		{
			name:       "generate recovery code",
			err:        NewErrGenRecoveryCode(internalErr),
			wantCode:   "auth_generate_recovery_code_failed",
			wantMsg:    "internal error: cannot generate and save the recovery code",
			wantStatus: http.StatusInternalServerError,
			wantUnwrap: internalErr,
		},
		{
			name:       "set verification",
			err:        NewErrSetVerification(internalErr),
			wantCode:   "auth_set_verification_failed",
			wantMsg:    "cannot update the verification state",
			wantStatus: http.StatusInternalServerError,
			wantUnwrap: internalErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertErrorShape(t, tt.err, tt.wantCode, tt.wantMsg, tt.wantStatus)
			assertErrorExtractors(t, tt.err)
			assertErrorIs(t, tt.err, tt.wantCode, tt.wantStatus)
			if !errors.Is(tt.err, internalErr) {
				t.Fatalf("errors.Is(err, internalErr) = false, want true")
			}
			if !errors.Is(tt.err.Unwrap(), tt.wantUnwrap) {
				t.Fatalf("Unwrap() = %v, want %v", tt.err.Unwrap(), tt.wantUnwrap)
			}
		})
	}
}
