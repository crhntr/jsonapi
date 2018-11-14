package jsonapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crhntr/jsonapi"
)

func TestHandle_ServeHTTP(t *testing.T) {
	t.Run("When Responding", func(t *testing.T) {
		// Setup
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Accept", jsonapi.ContentType)
		req.Header.Set("Content-Type", jsonapi.ContentType)

		res := httptest.NewRecorder()

		var mux jsonapi.Mux

		// Run
		mux.ServeHTTP(res, req)

		// Test Expectaions
		result := res.Result()

		if err := json.NewDecoder(result.Body).Decode(&struct{}{}); err != nil {
			t.Error("it should return an object")
		}
	})

	t.Run("When Correct Accept and Content Type Headers are Set", func(t *testing.T) {
		// Setup
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Accept", jsonapi.ContentType)
		req.Header.Set("Content-Type", jsonapi.ContentType)

		res := httptest.NewRecorder()

		var mux jsonapi.Mux

		// Run
		mux.ServeHTTP(res, req)

		// Test Expectaions
		result := res.Result()

		if result.StatusCode != http.StatusOK {
			t.Error("It should not return an error status code")
			t.Logf("instead got %q", result.StatusCode)
		}
	})

	t.Run("When Accept Request Header is Empty", func(t *testing.T) {
		// Setup
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Error(err)
		}

		res := httptest.NewRecorder()

		var mux jsonapi.Mux

		// Run
		mux.ServeHTTP(res, req)

		// Test Expectaions
		result := res.Result()
		if result.StatusCode != http.StatusNotAcceptable {
			t.Error("It should return http status not acceptable")
			t.Logf("instead got %q", result.StatusCode)
		}
	})

	t.Run("When Accept Request Header is Not Correct", func(t *testing.T) {
		// Setup
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Error(err)
		}

		req.Header.Set("Accept", "application/json")

		res := httptest.NewRecorder()

		var mux jsonapi.Mux

		// Run
		mux.ServeHTTP(res, req)

		// Test Expectaions
		result := res.Result()
		if result.StatusCode != http.StatusNotAcceptable {
			t.Error("It should return http status not acceptable")
		}
	})

	t.Run("When Content-Type Request Header is Empty", func(t *testing.T) {
		// Setup
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Accept", jsonapi.ContentType)

		res := httptest.NewRecorder()

		var mux jsonapi.Mux

		// Run
		mux.ServeHTTP(res, req)

		// Test Expectaions
		result := res.Result()
		if result.StatusCode != http.StatusUnsupportedMediaType {
			t.Error("It should return http status not unsupported media type")
		}
	})

	t.Run("When Content-Type Request Header is Not Correct", func(t *testing.T) {
		// Setup
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Accept", jsonapi.ContentType)
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()

		var mux jsonapi.Mux

		// Run
		mux.ServeHTTP(res, req)

		// Test Expectaions
		result := res.Result()
		if result.StatusCode != http.StatusUnsupportedMediaType {
			t.Error("It should return http status not unsupported media type")
		}
	})
}

func TestHandle_Resource(t *testing.T) {

}
