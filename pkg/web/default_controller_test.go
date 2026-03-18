package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

func TestDefaultControllerParamID(t *testing.T) {
	ctrl := DefaultController{}

	tests := []struct {
		name       string
		pathID     string
		wantErr    bool
		wantDetail string
	}{
		{
			name:       "missing id",
			wantErr:    true,
			wantDetail: "id is required",
		},
		{
			name:       "invalid id",
			pathID:     "invalid-uuid",
			wantErr:    true,
			wantDetail: "invalid id format:",
		},
		{
			name:   "valid id",
			pathID: uuid.NewString(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.pathID != "" {
				req.SetPathValue("id", tt.pathID)
			}

			id, err := ctrl.ParamID(req)
			if tt.wantErr {
				if err == nil {
					t.Fatal("ParamID() error = nil, want error")
				}

				httpErr, ok := err.(*fuego.HTTPError)
				if !ok {
					t.Fatalf("ParamID() error type = %T, want *fuego.HTTPError", err)
				}
				if httpErr.Status != http.StatusBadRequest {
					t.Fatalf("Status = %d, want %d", httpErr.Status, http.StatusBadRequest)
				}
				if !strings.Contains(httpErr.Detail, tt.wantDetail) {
					t.Fatalf("Detail = %q, want substring %q", httpErr.Detail, tt.wantDetail)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParamID() unexpected error = %v", err)
			}
			if id.String() != tt.pathID {
				t.Fatalf("ID = %q, want %q", id.String(), tt.pathID)
			}
		})
	}
}

func TestDefaultControllerParseHexID(t *testing.T) {
	ctrl := DefaultController{}
	wantID := uuid.NewString()

	got, err := ctrl.ParseHexID(wantID)
	if err != nil {
		t.Fatalf("ParseHexID() unexpected error = %v", err)
	}
	if got.String() != wantID {
		t.Fatalf("ID = %q, want %q", got.String(), wantID)
	}
}
