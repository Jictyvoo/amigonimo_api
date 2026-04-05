package smtp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
)

func withBoundary(contentType, boundary string) string {
	return contentType + `; boundary="` + boundary + `"`
}

// emailParts is the internal representation of a mail message.
type emailParts struct {
	sender          string
	to              []string
	subject         string
	body            string
	bodyHTML        string
	attachments     []attachment // regular file attachments (multipart/mixed)
	htmlAttachments []attachment // inline resources embedded in HTML (multipart/related)
}

// isMixed reports whether the message carries regular file attachments.
func (e emailParts) isMixed() bool { return len(e.attachments) > 0 }

// isAlternative reports whether the message provides both plain-text and HTML bodies.
func (e emailParts) isAlternative() bool { return len(e.body) > 0 && len(e.bodyHTML) > 0 }

// isRelated reports whether the HTML body references inline resources.
func (e emailParts) isRelated() bool { return len(e.bodyHTML) > 0 && len(e.htmlAttachments) > 0 }

// writeHeader serialises the RFC-2822 envelope headers via headerToBytes so that
// non-ASCII values (e.g. international subjects) are properly Q-encoded.
func (e emailParts) writeHeader(buf *bytes.Buffer) {
	header := textproto.MIMEHeader{}
	header.Set("From", e.sender)
	if len(e.to) > 0 {
		header.Set("To", e.to[0])
	}
	header.Set("Subject", e.subject)
	headerToBytes(buf, header)
}

// writePart writes a single body section as quoted-printable into mw (if non-nil)
// or directly into buf.
func (e emailParts) writePart(
	_ context.Context,
	mw *multipart.Writer,
	buf io.Writer,
	body, mediaType string,
) (err error) {
	pw := buf
	if mw != nil {
		header := textproto.MIMEHeader{
			mimeContentType:             {mediaType},
			mimeContentTransferEncoding: {mimeQP},
		}
		if pw, err = mw.CreatePart(header); err != nil {
			return err
		}
	}
	qw := quotedprintable.NewWriter(pw)
	defer func(qw *quotedprintable.Writer) {
		closeErr := qw.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}(qw)
	if _, err = io.WriteString(qw, body); err != nil {
		return err
	}
	return nil
}

// buildBody writes the text/HTML body sections into the multipart writer mw.
// mw may be nil when the message has no multipart structure (single plain-text body).
func (e emailParts) buildBody(ctx context.Context, mw *multipart.Writer, buf io.Writer) error {
	hasBoth := e.isAlternative()
	isMixed := e.isMixed()

	// When mixed AND has both bodies, wrap them in a nested multipart/alternative
	// part so clients can choose plain vs HTML.
	// Guard against nil mw: isMixed implies mw was created in buildMail, but the
	// nil guard keeps the static analyser (and future callers) safe.
	altWriter := mw
	if mw != nil && isMixed && hasBoth {
		altWriter = multipart.NewWriter(buf)
		defer altWriter.Close() //nolint:errcheck
		header := textproto.MIMEHeader{
			mimeContentType: {withBoundary(mimeMultiAlt, altWriter.Boundary())},
		}
		if _, err := mw.CreatePart(header); err != nil {
			return err
		}
	}

	if len(e.body) > 0 {
		bodyMW := (*multipart.Writer)(nil)
		if isMixed || hasBoth {
			bodyMW = altWriter
		}
		if err := e.writePart(ctx, bodyMW, buf, e.body, mimeTextPlain); err != nil {
			return err
		}
	}

	if len(e.bodyHTML) > 0 {
		return e.buildHTML(ctx, mw, altWriter, buf)
	}
	return nil
}

// buildHTML writes the HTML body and any inline (related) attachments.
func (e emailParts) buildHTML(
	ctx context.Context,
	mw, subMW *multipart.Writer,
	buf io.Writer,
) error {
	isMixed := e.isMixed()
	isAlt := e.isAlternative()
	hasInline := len(e.htmlAttachments) > 0

	// Default: write HTML directly into subMW.
	targetMW := subMW
	relatedMW := (*multipart.Writer)(nil)

	if hasInline && (isMixed || isAlt) {
		// Nested multipart/related inside the mixed/alt writer.
		relatedMW = multipart.NewWriter(buf)
		defer relatedMW.Close() //nolint:errcheck
		header := textproto.MIMEHeader{
			mimeContentType: {withBoundary(mimeMultiRelated, relatedMW.Boundary())},
		}
		if _, err := subMW.CreatePart(header); err != nil {
			return err
		}
		targetMW = relatedMW
	} else if e.isRelated() {
		// Top-level related (no mixed/alt wrapper): reuse mw.
		relatedMW = mw
		targetMW = mw
	}

	htmlBodyMW := (*multipart.Writer)(nil)
	if isMixed || isAlt || e.isRelated() {
		htmlBodyMW = targetMW
	}
	if err := e.writePart(ctx, htmlBodyMW, buf, e.bodyHTML, mimeTextHTML); err != nil {
		return err
	}
	return e.writeInlineAttachments(relatedMW)
}

// writeInlineAttachments appends HTML-related (inline) attachments to relatedMW.
func (e emailParts) writeInlineAttachments(relatedMW *multipart.Writer) error {
	for i := range e.htmlAttachments {
		e.htmlAttachments[i].isHTMLRelated = true
		if err := e.htmlAttachments[i].build(relatedMW); err != nil {
			return err
		}
	}
	return nil
}

// buildMail serialises emailParts into a raw RFC-2822 byte slice.
func buildMail(ctx context.Context, e *emailParts) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 4096+len(e.body)+len(e.bodyHTML)))

	var mw *multipart.Writer
	if e.isMixed() || e.isAlternative() || e.isRelated() {
		mw = multipart.NewWriter(buf)
		defer mw.Close() //nolint:errcheck
	}

	// Envelope headers (From, To, Subject).
	e.writeHeader(buf)

	// Top-level MIME content headers: pick the correct multipart subtype.
	switch {
	case mw != nil:
		topType := mimeMultiMixed
		switch {
		case !e.isMixed() && e.isAlternative():
			topType = mimeMultiAlt
		case !e.isMixed() && !e.isAlternative() && e.isRelated():
			topType = mimeMultiRelated
		}
		headerToBytes(buf, textproto.MIMEHeader{
			mimeContentType: {withBoundary(topType, mw.Boundary())},
		})
	case len(e.bodyHTML) > 0:
		headerToBytes(
			buf, textproto.MIMEHeader{
				mimeContentType:             {mimeTextHTML},
				mimeContentTransferEncoding: {mimeQP},
			},
		)
	default:
		headerToBytes(
			buf, textproto.MIMEHeader{
				mimeContentType:             {mimeTextPlain},
				mimeContentTransferEncoding: {mimeQP},
			},
		)
	}
	buf.WriteString("\r\n") // blank line between headers and body

	if err := e.buildBody(ctx, mw, buf); err != nil {
		return nil, err
	}

	for i := range e.attachments {
		e.attachments[i].isHTMLRelated = false
		if err := e.attachments[i].build(mw); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// headerToBytes writes MIME header fields to a StringWriter.
// Non-structural headers (everything except Content-Type and Content-Disposition)
// are Q-encoded so that non-ASCII values (e.g. international subjects) are safe.
func headerToBytes(w io.StringWriter, header textproto.MIMEHeader) {
	for field, values := range header {
		for _, v := range values {
			_, _ = w.WriteString(field + ": ")
			switch field {
			case mimeContentType, mimeContentDisposition:
				_, _ = w.WriteString(v)
			default:
				_, _ = w.WriteString(mime.QEncoding.Encode("UTF-8", v))
			}
			_, _ = w.WriteString("\r\n")
		}
	}
}
