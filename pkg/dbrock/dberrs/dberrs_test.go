package dberrs

import (
	"errors"
	"net/http"
	"testing"
)

type mockTestError struct {
	Message string
}

func (e mockTestError) Error() string {
	return e.Message
}

func TestBaseErrorCheckerIs_Comprehensive(t *testing.T) {
	type customTestErrorChecker struct {
		baseErrorWrapper[customTestErrorChecker, *customTestErrorChecker]
		mockTestError
	}

	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{
			name:   "same raw error type",
			err:    customTestErrorChecker{mockTestError: mockTestError{Message: "hello world"}},
			target: customTestErrorChecker{},
			want:   true,
		},
		{
			name:   "same pointer error type",
			err:    customTestErrorChecker{mockTestError: mockTestError{Message: "hello world"}},
			target: &customTestErrorChecker{},
			want:   true,
		},
		{
			name: "same error as pointer type",
			err: customTestErrorChecker{
				mockTestError: mockTestError{Message: "not nil-pointer"},
			},
			target: &customTestErrorChecker{},
			want:   true,
		},
		{
			name:   "check pointer generated with raw type",
			err:    customTestErrorChecker{mockTestError: mockTestError{Message: "pointer"}},
			target: customTestErrorChecker{},
			want:   true,
		},
		{
			name:   "different error type",
			err:    customTestErrorChecker{},
			target: errors.New("some other error"),
			want:   false,
		},
		{
			name:   "nil error",
			err:    nil,
			target: customTestErrorChecker{},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.err, tt.target); got != tt.want {
				t.Errorf("errors.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseErrorWrapperUnwrap(t *testing.T) {
	type wrapErrFull struct {
		baseErrorWrapper[wrapErrFull, *wrapErrFull]
		mockTestError
	}

	inner := errors.New("inner")
	e := wrapErrFull{}
	e.Err = inner

	if !errors.Is(e.Unwrap(), inner) {
		t.Errorf("Unwrap() = %v, want %v", e.Unwrap(), inner)
	}

	e2 := wrapErrFull{}
	if got := e2.Unwrap(); got != nil {
		t.Errorf("Unwrap() on nil inner = %v, want nil", got)
	}
}

func TestBaseErrorWrapperStatusCode(t *testing.T) {
	type wrapErrFull struct {
		baseErrorWrapper[wrapErrFull, *wrapErrFull]
		mockTestError
	}
	e := wrapErrFull{}
	if got := e.StatusCode(); got != http.StatusInternalServerError {
		t.Errorf("StatusCode() = %d, want %d", got, http.StatusInternalServerError)
	}
}

func TestBaseErrorWrapperDetailMsg(t *testing.T) {
	type wrapErrFull struct {
		baseErrorWrapper[wrapErrFull, *wrapErrFull]
		mockTestError
	}

	tests := []struct {
		name  string
		inner error
		want  string
	}{
		{"nil inner returns unknown error", nil, "unknown error"},
		{"non-nil inner returns its message", errors.New("boom"), "boom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := wrapErrFull{}
			e.Err = tt.inner
			if got := e.DetailMsg(); got != tt.want {
				t.Errorf("DetailMsg() = %q, want %q", got, tt.want)
			}
		})
	}
}
