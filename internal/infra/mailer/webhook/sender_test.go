package webhook_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer"
	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer/webhook"
)

// captured holds what the test server received.
type captured struct {
	method      string
	contentType string
	authHeader  string
	body        []byte
}

// newTestServer starts an httptest server that captures the request and
// responds with statusCode.
func newTestServer(t *testing.T, statusCode int) (*httptest.Server, *captured) {
	t.Helper()
	cap := &captured{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cap.method = r.Method
		cap.contentType = r.Header.Get("Content-Type")
		cap.authHeader = r.Header.Get("Authorization")
		cap.body, _ = io.ReadAll(r.Body)
		w.WriteHeader(statusCode)
	}))
	t.Cleanup(srv.Close)
	return srv, cap
}

// ── Send ─────────────────────────────────────────────────────────────────────

func TestWebhookSend(t *testing.T) {
	tests := []struct {
		name       string
		apiKey     string
		statusCode int
		msg        mailer.EmailMessage
		wantErr    bool
		errContain string
	}{
		{
			name:       "200 OK — no error",
			statusCode: http.StatusOK,
			msg:        mailer.EmailMessage{To: []string{"u@example.com"}, Subject: "Hi"},
		},
		{
			name:       "201 Created — no error",
			statusCode: http.StatusCreated,
			msg:        mailer.EmailMessage{To: []string{"u@example.com"}, Subject: "Hi"},
		},
		{
			name:       "400 Bad Request — error",
			statusCode: http.StatusBadRequest,
			msg:        mailer.EmailMessage{To: []string{"u@example.com"}, Subject: "Hi"},
			wantErr:    true,
			errContain: "400",
		},
		{
			name:       "500 Internal Server Error — error",
			statusCode: http.StatusInternalServerError,
			msg:        mailer.EmailMessage{To: []string{"u@example.com"}, Subject: "Hi"},
			wantErr:    true,
			errContain: "500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, _ := newTestServer(t, tt.statusCode)
			s := webhook.New(webhook.Config{URL: srv.URL, APIKey: tt.apiKey})

			err := s.Send(context.Background(), tt.msg)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.errContain != "" && err != nil && !strings.Contains(err.Error(), tt.errContain) {
				t.Errorf("error %q does not contain %q", err.Error(), tt.errContain)
			}
		})
	}
}

func TestWebhookSendNetworkError(t *testing.T) {
	s := webhook.New(webhook.Config{URL: "http://127.0.0.1:1"}) // nothing listening
	err := s.Send(
		context.Background(),
		mailer.EmailMessage{To: []string{"u@example.com"}, Subject: "Hi"},
	)
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// ── Request payload ───────────────────────────────────────────────────────────

func TestWebhookPayload(t *testing.T) {
	tests := []struct {
		name    string
		msg     mailer.EmailMessage
		apiKey  string
		checkFn func(t *testing.T, cap *captured)
	}{
		{
			name: "payload fields mapped correctly",
			msg: mailer.EmailMessage{
				To:       []string{"a@example.com", "b@example.com"},
				Subject:  "Test subject",
				Body:     "Plain text",
				BodyHTML: "<p>HTML</p>",
			},
			checkFn: func(t *testing.T, cap *captured) {
				var payload map[string]any
				if err := json.Unmarshal(cap.body, &payload); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				to, _ := payload["to"].([]any)
				if len(to) != 2 {
					t.Errorf("to length = %d, want 2", len(to))
				}
				if payload["subject"] != "Test subject" {
					t.Errorf("subject = %v", payload["subject"])
				}
				if payload["body"] != "Plain text" {
					t.Errorf("body = %v", payload["body"])
				}
				if payload["body_html"] != "<p>HTML</p>" {
					t.Errorf("body_html = %v", payload["body_html"])
				}
			},
		},
		{
			name: "Content-Type header is application/json",
			msg:  mailer.EmailMessage{To: []string{"u@example.com"}, Subject: "Hi"},
			checkFn: func(t *testing.T, cap *captured) {
				if cap.contentType != "application/json" {
					t.Errorf("Content-Type = %q, want application/json", cap.contentType)
				}
			},
		},
		{
			name:   "Authorization header set when APIKey provided",
			msg:    mailer.EmailMessage{To: []string{"u@example.com"}, Subject: "Hi"},
			apiKey: "secret-key",
			checkFn: func(t *testing.T, cap *captured) {
				want := "Bearer secret-key"
				if cap.authHeader != want {
					t.Errorf("Authorization = %q, want %q", cap.authHeader, want)
				}
			},
		},
		{
			name: "Authorization header absent when APIKey empty",
			msg:  mailer.EmailMessage{To: []string{"u@example.com"}, Subject: "Hi"},
			checkFn: func(t *testing.T, cap *captured) {
				if cap.authHeader != "" {
					t.Errorf("Authorization = %q, want empty", cap.authHeader)
				}
			},
		},
		{
			name: "HTTP method is POST",
			msg:  mailer.EmailMessage{To: []string{"u@example.com"}, Subject: "Hi"},
			checkFn: func(t *testing.T, cap *captured) {
				if cap.method != http.MethodPost {
					t.Errorf("method = %q, want POST", cap.method)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, cap := newTestServer(t, http.StatusOK)
			s := webhook.New(webhook.Config{URL: srv.URL, APIKey: tt.apiKey})

			if err := s.Send(context.Background(), tt.msg); err != nil {
				t.Fatalf("Send() unexpected error: %v", err)
			}
			tt.checkFn(t, cap)
		})
	}
}

func TestWebhookAttachmentPayload(t *testing.T) {
	data := []byte("hello attachment")
	msg := mailer.EmailMessage{
		To:      []string{"u@example.com"},
		Subject: "with attachment",
		Body:    "see attached",
		Attachments: []mailer.Attachment{
			{Filename: "file.txt", ContentType: "text/plain", Data: data, Inline: false},
			{Filename: "img.png", ContentType: "image/png", Data: []byte("pngdata"), Inline: true},
		},
	}

	srv, cap := newTestServer(t, http.StatusOK)
	s := webhook.New(webhook.Config{URL: srv.URL})

	if err := s.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send() unexpected error: %v", err)
	}

	type attPayload struct {
		Filename    string `json:"filename"`
		ContentType string `json:"content_type"`
		Data        string `json:"data"`
		Inline      bool   `json:"inline"`
	}
	var payload struct {
		Attachments []attPayload `json:"attachments"`
	}
	if err := json.Unmarshal(cap.body, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(payload.Attachments) != 2 {
		t.Fatalf("attachments length = %d, want 2", len(payload.Attachments))
	}

	first := payload.Attachments[0]
	if first.Filename != "file.txt" {
		t.Errorf("attachment[0].filename = %q, want file.txt", first.Filename)
	}
	decoded, err := base64.StdEncoding.DecodeString(first.Data)
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}
	if string(decoded) != string(data) {
		t.Errorf("attachment data = %q, want %q", decoded, data)
	}
	if payload.Attachments[1].Inline != true {
		t.Error("inline attachment should have inline=true")
	}
}

func TestWebhookContextCancellation(t *testing.T) {
	// A server that hangs; the cancelled context should abort the request.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	s := webhook.New(webhook.Config{URL: srv.URL})
	err := s.Send(ctx, mailer.EmailMessage{To: []string{"u@example.com"}, Subject: "Hi"})
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
	if !errors.Is(err, context.Canceled) && !strings.Contains(err.Error(), "context") {
		t.Errorf("expected context-related error, got: %v", err)
	}
}
