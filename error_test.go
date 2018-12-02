package jsonapi_test

import (
	"net/http"
	"testing"

	"github.com/crhntr/jsonapi"
)

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
