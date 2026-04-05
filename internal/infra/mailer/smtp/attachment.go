package smtp

import (
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/textproto"
)

// attachment represents a file or inline resource attached to an email.
type attachment struct {
	contentType   string
	filename      string
	headers       textproto.MIMEHeader
	data          []byte
	isHTMLRelated bool
}

func (a *attachment) normalizeHeaders() {
	if a.headers == nil {
		a.headers = textproto.MIMEHeader{}
	}

	ct := mimeOctetStream
	if a.contentType != "" {
		ct = a.contentType
	}
	a.headers.Set(mimeContentType, ct)

	if a.headers.Get(mimeContentDisposition) == "" {
		disposition := mimeAttachment
		if a.isHTMLRelated {
			disposition = mimeInline
		}
		a.headers.Set(
			mimeContentDisposition,
			fmt.Sprintf("%s;\r\n filename=%q", disposition, a.filename),
		)
	}
	if a.headers.Get(mimeContentID) == "" {
		a.headers.Set(mimeContentID, fmt.Sprintf("<%s>", a.filename))
	}
	if a.headers.Get(mimeContentTransferEncoding) == "" {
		a.headers.Set(mimeContentTransferEncoding, mimeBase64)
	}
}

func (a *attachment) build(mw *multipart.Writer) error {
	a.normalizeHeaders()
	pw, err := mw.CreatePart(a.headers)
	if err != nil {
		return err
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(a.data)))
	base64.StdEncoding.Encode(encoded, a.data)
	_, err = pw.Write(encoded)
	return err
}
