package fixtures

// ErrorDetail is a minimal error response type used in failure assertions.
// It only matches the "detail" field of the server error response,
// so tests don't need to specify the full fuego.HTTPError shape.
type ErrorDetail struct {
	Detail string `json:"detail"`
}
