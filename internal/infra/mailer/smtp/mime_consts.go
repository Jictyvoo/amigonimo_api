package smtp

// MIME constants.
const (
	mimeContentType             = "Content-Type"
	mimeContentTransferEncoding = "Content-Transfer-Encoding"
	mimeContentDisposition      = "Content-Disposition"
	mimeContentID               = "Content-ID"

	mimeTextPlain   = `text/plain; charset="utf-8"`
	mimeTextHTML    = `text/html; charset="utf-8"`
	mimeQP          = "quoted-printable"
	mimeBase64      = "base64"
	mimeOctetStream = "application/octet-stream"
	mimeAttachment  = "attachment"
	mimeInline      = "inline"
)

// Multipart subtypes:
//
//	mixed - body + regular file attachments (PDFs)
//	alternative - plain-text AND HTML versions of the same message
//	related - HTML body and inline resources (embedded images)
const (
	mimeMultiMixed   = "multipart/mixed"
	mimeMultiAlt     = "multipart/alternative"
	mimeMultiRelated = "multipart/related"
)
