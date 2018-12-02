package jsonapi_test

import (
	"net/http"
	"testing"

	"github.com/crhntr/jsonapi"
)

func TestError_HTTPStatus(t *testing.T) {
	t.Run("when it has a zero value", func(t *testing.T) {
		var error jsonapi.Error
		if status := error.HTTPStatus(); status != http.StatusInternalServerError {
			t.Error("it should return StatusInternalServerError")
			t.Log(status)
		}
	})

	t.Run("when it has a non zero status", func(t *testing.T) {
		error := jsonapi.Error{Status: http.StatusTeapot}
		if status := error.HTTPStatus(); status != http.StatusTeapot {
			t.Error("it should return StatusTeapot")
			t.Log(status)
		}
	})
}

func TestError_Error(t *testing.T) {
	t.Run("when it has a zero value", func(t *testing.T) {
		var error jsonapi.Error
		if msg := error.Error(); msg != "{}" {
			t.Error("it should return an empty object")
			t.Log(msg)
		}
	})

	t.Run("when it encounteres a Marshalling error", func(t *testing.T) {
		var error jsonapi.Error
		error.About = &jsonapi.Link{}
		expectedMessage := "could not encode error: json: error calling MarshalJSON for type *jsonapi.Link: a link must have a string or object value"
		if msg := error.Error(); msg != expectedMessage {
			t.Error("it should an error message including the json.Marshaling error")
			t.Log(msg)
		}
	})
}

func TestErrorsPolicy(t *testing.T) {
	t.Run("when an empty list of errors are passed", func(t *testing.T) {
		if status := jsonapi.ErrorsPolicy(nil); status != http.StatusOK {
			t.Error("it should return Status OK")
			t.Log(status)
		}
	})

	t.Run("when a single error is passed", func(t *testing.T) {
		errors := []jsonapi.Error{{Status: http.StatusBadRequest}}
		if status := jsonapi.ErrorsPolicy(errors); status != http.StatusBadRequest {
			t.Error("it should return the error status")
			t.Log(status)
		}
	})

	t.Run("when a multiple errors with bad request statuses", func(t *testing.T) {
		errors := []jsonapi.Error{{Status: http.StatusTeapot}, {Status: http.StatusUnauthorized}}
		if status := jsonapi.ErrorsPolicy(errors); status != http.StatusBadRequest {
			t.Error("it should return the most general status")
			t.Log(status)
		}
	})

	t.Run("when a multiple errors with server error statuses", func(t *testing.T) {
		errors := []jsonapi.Error{{Status: http.StatusServiceUnavailable}, {Status: http.StatusInsufficientStorage}}
		if status := jsonapi.ErrorsPolicy(errors); status != http.StatusInternalServerError {
			t.Error("it should return the most general status")
			t.Log(status)
		}
	})

	t.Run("when a multiple errors with both 4XX and 5XX", func(t *testing.T) {
		errors := []jsonapi.Error{
			{Status: http.StatusServiceUnavailable}, {Status: http.StatusInsufficientStorage},
			{Status: http.StatusTeapot}, {Status: http.StatusUnauthorized},
		}
		if status := jsonapi.ErrorsPolicy(errors); status != http.StatusInternalServerError {
			t.Error("it should return the most general highest status")
			t.Log(status)
		}
	})
}
