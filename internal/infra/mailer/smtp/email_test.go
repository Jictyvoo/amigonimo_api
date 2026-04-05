package smtp

import (
	"bytes"
	"context"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
	"testing"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func newParts(body, bodyHTML string, atts, htmlAtts []attachment) emailParts {
	return emailParts{
		sender:          "Sender <sender@example.com>",
		to:              []string{"to@example.com"},
		subject:         "Test subject",
		body:            body,
		bodyHTML:        bodyHTML,
		attachments:     atts,
		htmlAttachments: htmlAtts,
	}
}

func mustBuildMail(t *testing.T, e emailParts) []byte {
	t.Helper()
	raw, err := buildMail(context.Background(), &e)
	if err != nil {
		t.Fatalf("buildMail error: %v", err)
	}
	return raw
}

func parseHeaders(t *testing.T, raw []byte) mail.Header {
	t.Helper()
	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("mail.ReadMessage error: %v", err)
	}
	return msg.Header
}

func assertHeaderContains(t *testing.T, header mail.Header, key, want string) {
	t.Helper()
	v := header.Get(key)
	if !strings.Contains(v, want) {
		t.Errorf("header %q = %q; want it to contain %q", key, v, want)
	}
}

func assertRawContains(t *testing.T, raw []byte, substr string) {
	t.Helper()
	if !bytes.Contains(raw, []byte(substr)) {
		t.Errorf("raw email does not contain %q", substr)
	}
}

// ── withBoundary ─────────────────────────────────────────────────────────────

func TestWithBoundary(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		boundary    string
		want        string
	}{
		{
			name:        "mixed",
			contentType: mimeMultiMixed,
			boundary:    "abc123",
			want:        `multipart/mixed; boundary="abc123"`,
		},
		{
			name:        "alternative",
			contentType: mimeMultiAlt,
			boundary:    "xyz",
			want:        `multipart/alternative; boundary="xyz"`,
		},
		{
			name:        "related",
			contentType: mimeMultiRelated,
			boundary:    "def",
			want:        `multipart/related; boundary="def"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := withBoundary(tt.contentType, tt.boundary)
			if got != tt.want {
				t.Errorf("withBoundary() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ── emailParts predicates ─────────────────────────────────────────────────────

func TestEmailPartsPredicates(t *testing.T) {
	att := attachment{filename: "f.txt", data: []byte("x")}

	tests := []struct {
		name        string
		parts       emailParts
		wantMixed   bool
		wantAlt     bool
		wantRelated bool
	}{
		{
			name:  "plain text only",
			parts: newParts("hello", "", nil, nil),
		},
		{
			name:    "html only — not alternative",
			parts:   newParts("", "<p>hi</p>", nil, nil),
			wantAlt: false,
		},
		{
			name:    "text + html → alternative",
			parts:   newParts("hi", "<p>hi</p>", nil, nil),
			wantAlt: true,
		},
		{
			name:      "file attachment → mixed",
			parts:     newParts("hi", "", []attachment{att}, nil),
			wantMixed: true,
		},
		{
			name:        "html + inline → related",
			parts:       newParts("", "<p>hi</p>", nil, []attachment{att}),
			wantRelated: true,
		},
		{
			name:        "all three",
			parts:       newParts("hi", "<p>hi</p>", []attachment{att}, []attachment{att}),
			wantMixed:   true,
			wantAlt:     true,
			wantRelated: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.parts.isMixed(); got != tt.wantMixed {
				t.Errorf("isMixed() = %v, want %v", got, tt.wantMixed)
			}
			if got := tt.parts.isAlternative(); got != tt.wantAlt {
				t.Errorf("isAlternative() = %v, want %v", got, tt.wantAlt)
			}
			if got := tt.parts.isRelated(); got != tt.wantRelated {
				t.Errorf("isRelated() = %v, want %v", got, tt.wantRelated)
			}
		})
	}
}

// ── headerToBytes ─────────────────────────────────────────────────────────────

func TestHeaderToBytes(t *testing.T) {
	tests := []struct {
		name    string
		header  map[string][]string
		wantSub string // expected substring in output
	}{
		{
			name:    "content-type is not encoded",
			header:  map[string][]string{mimeContentType: {"text/plain"}},
			wantSub: "Content-Type: text/plain",
		},
		{
			name:    "subject is q-encoded",
			header:  map[string][]string{"Subject": {"Hello"}},
			wantSub: "Subject: ",
		},
		{
			name:    "content-disposition is not encoded",
			header:  map[string][]string{mimeContentDisposition: {`attachment; filename="f.txt"`}},
			wantSub: `attachment; filename="f.txt"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &strings.Builder{}
			headerToBytes(buf, tt.header)
			if !strings.Contains(buf.String(), tt.wantSub) {
				t.Errorf("headerToBytes output %q does not contain %q", buf.String(), tt.wantSub)
			}
		})
	}
}

// ── buildMail ─────────────────────────────────────────────────────────────────

func TestBuildMail(t *testing.T) {
	att := attachment{
		filename:    "report.pdf",
		contentType: "application/pdf",
		data:        []byte("PDFDATA"),
	}
	inlineAtt := attachment{filename: "img.png", contentType: "image/png", data: []byte("PNGDATA")}

	tests := []struct {
		name           string
		parts          emailParts
		wantCTContains string // substring expected in Content-Type header
		wantInRaw      []string
	}{
		{
			name:           "plain text only",
			parts:          newParts("Hello plain", "", nil, nil),
			wantCTContains: "text/plain",
			wantInRaw:      []string{"Hello plain", "Subject: "},
		},
		{
			name:           "html only",
			parts:          newParts("", "<b>Hello</b>", nil, nil),
			wantCTContains: "text/html",
			wantInRaw:      []string{"<b>Hello</b>"},
		},
		{
			name:           "text and html → alternative",
			parts:          newParts("Plain body", "<p>HTML body</p>", nil, nil),
			wantCTContains: "multipart/alternative",
			wantInRaw:      []string{"Plain body", "<p>HTML body</p>"},
		},
		{
			name:           "text with file attachment → mixed",
			parts:          newParts("See attached", "", []attachment{att}, nil),
			wantCTContains: "multipart/mixed",
			wantInRaw:      []string{"See attached", "report.pdf"},
		},
		{
			name:           "html with inline attachment → related",
			parts:          newParts("", "<img src='cid:img.png'>", nil, []attachment{inlineAtt}),
			wantCTContains: "multipart",
			wantInRaw:      []string{"img.png"},
		},
		{
			name:           "envelope headers present",
			parts:          newParts("body", "", nil, nil),
			wantCTContains: "text/plain",
			wantInRaw:      []string{"sender@example.com", "to@example.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := mustBuildMail(t, tt.parts)
			headers := parseHeaders(t, raw)

			assertHeaderContains(t, headers, "Content-Type", tt.wantCTContains)
			for _, sub := range tt.wantInRaw {
				assertRawContains(t, raw, sub)
			}
		})
	}
}

func TestBuildMailMultipartStructure(t *testing.T) {
	att := attachment{
		filename:    "file.txt",
		contentType: "text/plain",
		data:        []byte("attachment content"),
	}

	raw := mustBuildMail(t, newParts("text body", "<p>html body</p>", []attachment{att}, nil))

	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("ReadMessage: %v", err)
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		t.Fatalf("ParseMediaType: %v", err)
	}
	if !strings.HasPrefix(mediaType, "multipart/") {
		t.Fatalf("expected multipart, got %q", mediaType)
	}

	mr := multipart.NewReader(msg.Body, params["boundary"])
	partCount := 0
	for {
		p, readErr := mr.NextPart()
		if readErr != nil {
			break
		}
		p.Close()
		partCount++
	}
	if partCount == 0 {
		t.Error("expected at least one multipart part")
	}
}

// ── mime_consts sanity ────────────────────────────────────────────────────────

func TestMimeConsts(t *testing.T) {
	// Ensure no typos in the declared MIME type constants.
	types := []string{mimeTextPlain, mimeTextHTML, mimeMultiMixed, mimeMultiAlt, mimeMultiRelated}
	for _, ct := range types {
		t.Run(ct, func(t *testing.T) {
			// mime.ParseMediaType must not return an error.
			mediaType, _, err := mime.ParseMediaType(ct)
			if err != nil {
				t.Errorf("mime.ParseMediaType(%q) error: %v", ct, err)
			}
			if mediaType == "" {
				t.Errorf("mime.ParseMediaType(%q) returned empty media type", ct)
			}
		})
	}
}
