package reqrunner

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/stretchr/testify/assert"

	"github.com/jictyvoo/amigonimo_api/test/internal/runners"
)

// Validator defines how to validate an HTTP response.
type (
	respObjectSanitizer[O any] func(expected, actual *O) error
	Validator                  interface {
		Validate(rCtx runners.RunnerContext, resp *http.Response)
	}
)

// jsonBodyValidator defines expected results for a request body.
type jsonBodyValidator[O any] struct {
	Body      O
	Sanitizer respObjectSanitizer[O]
}

func (v jsonBodyValidator[O]) decodeBody(reader io.Reader) (O, error) {
	decoder := json.NewDecoder(reader)

	output := new(O)
	err := decoder.Decode(output)
	return *output, err
}

func (v jsonBodyValidator[O]) Validate(rCtx runners.RunnerContext, resp *http.Response) {
	// Unmarshal response
	actualBody, err := v.decodeBody(resp.Body)
	if err != nil {
		rCtx.Fatalf("Failed to decode response body: %v", err)
	}

	// Persist the decoded body to storage automatically
	runners.SaveOnCtx(rCtx, actualBody)

	// If Sanitizer is provided, use it
	expectedBody := v.Body
	if v.Sanitizer != nil {
		if err = v.Sanitizer(&expectedBody, &actualBody); err != nil {
			rCtx.Fatalf("Failed to sanitize response body: %v", err)
		}
	}

	assert.Equal(rCtx, expectedBody, actualBody, "response body mismatch")
}

// ExpectStatus adds a validator for the response status code.
func ExpectStatus(code int) Option {
	return func(r *HttpRunner) {
		r.validators = append(r.validators, statusValidator(code))
	}
}

type statusValidator int

func (v statusValidator) Validate(rCtx runners.RunnerContext, resp *http.Response) {
	if int(v) != resp.StatusCode {
		rCtx.Fatalf("Unexpected status code: got %v want %v", resp.StatusCode, int(v))
	}
}

// ExpectBody adds a validator for the response body.
func ExpectBody[O any](expected O, sanitizer ...respObjectSanitizer[O]) Option {
	return func(r *HttpRunner) {
		newValidator := jsonBodyValidator[O]{Body: expected}
		if len(sanitizer) > 0 {
			newValidator.Sanitizer = sanitizer[0]
		}

		r.validators = append(r.validators, newValidator)
	}
}

// ExtractToState adds a validator that extracts a value from a previously stored response
// and saves it to storage.
func ExtractToState[R any, V any](extractor func(resp R) V) Option {
	return func(r *HttpRunner) {
		r.validators = append(r.validators, &extractorValidator[R, V]{extractor: extractor})
	}
}

type extractorValidator[R any, V any] struct {
	extractor func(resp R) V
}

func (v *extractorValidator[R, V]) Validate(rCtx runners.RunnerContext, _ *http.Response) {
	respBody, ok := runners.LoadFromCtx[R](rCtx)
	if !ok {
		rCtx.Fatalf("Failed to load response of type %T from storage for extraction", respBody)
	}

	val := v.extractor(respBody)
	runners.SaveOnCtx(rCtx, val)
}
