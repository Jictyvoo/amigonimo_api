package webhook

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer"
)

var _ mailer.Sender = (*Sender)(nil)

const (
	headerContentType   = "Content-Type"
	headerAuthorization = "Authorization"
	contentTypeJSON     = "application/json"
	authBearer          = "Bearer "
)

type attachmentPayload struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type,omitempty"`
	Data        string `json:"data"` // base64-encoded
	Inline      bool   `json:"inline,omitempty"`
}

type emailPayload struct {
	To          []string            `json:"to"`
	Subject     string              `json:"subject"`
	Body        string              `json:"body,omitempty"`
	BodyHTML    string              `json:"body_html,omitempty"`
	Attachments []attachmentPayload `json:"attachments,omitempty"`
}

// Sender delivers email by POST-ing a JSON payload to a configured webhook URL.
type Sender struct {
	config Config
	client *http.Client
}

func New(cfg Config) *Sender {
	return &Sender{config: cfg, client: http.DefaultClient}
}

// buildRequest constructs the outgoing HTTP request from the email message.
func (s *Sender) buildRequest(ctx context.Context, msg mailer.EmailMessage) (*http.Request, error) {
	var atts []attachmentPayload
	for _, a := range msg.Attachments {
		atts = append(
			atts, attachmentPayload{
				Filename:    a.Filename,
				ContentType: a.ContentType,
				Data:        base64.StdEncoding.EncodeToString(a.Data),
				Inline:      a.Inline,
			},
		)
	}

	payload, err := json.Marshal(
		emailPayload{
			To:          msg.To,
			Subject:     msg.Subject,
			Body:        msg.Body,
			BodyHTML:    msg.BodyHTML,
			Attachments: atts,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("webhook mailer: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, s.config.URL, bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("webhook mailer: build request: %w", err)
	}
	req.Header.Set(headerContentType, contentTypeJSON)
	if s.config.APIKey != "" {
		req.Header.Set(headerAuthorization, authBearer+s.config.APIKey)
	}
	return req, nil
}

func (s *Sender) Send(ctx context.Context, msg mailer.EmailMessage) error {
	req, err := s.buildRequest(ctx, msg)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook mailer: send: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook mailer: unexpected status %d", resp.StatusCode)
	}
	return nil
}
