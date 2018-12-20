package jsonapi

import (
	"io"
	"net/http"
)

// NewRequest sets required jsonapi headers for requests to jsonapi servers.
// It is used internally for tests and examples.
func NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return req, err
	}
	req.Header.Set("Accept", ContentType)
	req.Header.Set("Content-Type", ContentType)
	return req, nil
}
