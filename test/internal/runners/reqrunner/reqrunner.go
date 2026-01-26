package reqrunner

import (
	"io"
	"net/http"

	"github.com/jictyvoo/amigonimo_api/test/internal/runners"
)

// RequestMaker defines how to construct an HTTP request.
type RequestMaker func(baseURL string, ctx runners.RunnerContext) (*http.Request, error)

// HttpRunner is the main HTTP test runner.
type HttpRunner struct {
	BaseURL     string
	Client      *http.Client
	makeRequest RequestMaker
	validators  []Validator
}

// Option is a functional option for configuring the HttpRunner.
type Option func(*HttpRunner)

// NewHttpRunner creates a new HttpRunner with the given options.
func NewHttpRunner(baseURL string, opts ...Option) *HttpRunner {
	r := &HttpRunner{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *HttpRunner) Run(rCtx runners.RunnerContext) error {
	if r.makeRequest == nil {
		rCtx.Fatal("No request configured for HttpRunner")
	}

	// Execute Request
	httpReq, err := r.makeRequest(r.BaseURL, rCtx)
	if err != nil {
		rCtx.Fatalf("Impossible to perform request: %v", err)
	}

	resp, reqErr := r.Client.Do(httpReq)
	if reqErr != nil {
		rCtx.Fatalf("Failed to perform request: %v", reqErr)
	}
	defer func(Body io.ReadCloser) {
		if closeErr := Body.Close(); closeErr != nil {
			rCtx.Errorf("Failed to close response body: %v", closeErr)
		}
	}(resp.Body)

	// Run validators
	for _, v := range r.validators {
		v.Validate(rCtx, resp)
	}

	return nil
}
