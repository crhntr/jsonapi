package jsonapi

import (
	"io"
	"net/http"
)

// NewRequest is a helper method to set required headers on requests
func NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return req, err
	}
	req.Header.Set("Accept", ContentType)
	req.Header.Set("Content-Type", ContentType)
	return req, nil
}
