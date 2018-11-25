package jsonapi_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crhntr/jsonapi"
)

func TestHandle_ServeHTTP_TopLevelAndContentNegotiation(t *testing.T) {
	t.Run("When Responding", func(t *testing.T) {
		// Setup
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Accept", jsonapi.ContentType)
		req.Header.Set("Content-Type", jsonapi.ContentType)

		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux

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

		var mux jsonapi.ServeMux

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

		var mux jsonapi.ServeMux

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

		var mux jsonapi.ServeMux

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

		var mux jsonapi.ServeMux

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

		var mux jsonapi.ServeMux

		// Run
		mux.ServeHTTP(res, req)

		// Test Expectaions
		result := res.Result()
		if result.StatusCode != http.StatusUnsupportedMediaType {
			t.Error("It should return http status not unsupported media type")
		}
	})
}

func TestHandle_ServeHTTP_RequestMux(t *testing.T) {
	t.Run("When GET empty resource collection", func(t *testing.T) {
		// Setup
		req, err := http.NewRequest(http.MethodGet, "/resource", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Accept", jsonapi.ContentType)
		req.Header.Set("Content-Type", jsonapi.ContentType)

		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux

		var calledHandler bool
		mux.HandleFetchCollection("resource", func(res jsonapi.FetchCollectionResponder, req *http.Request) {
			calledHandler = true
		})
		// Run
		mux.ServeHTTP(res, req)

		// Test Expectaions
		if !calledHandler {
			t.Error(`it should call handler`)
		}

		result := res.Result()

		bodyBuf, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Error("it should return valid json")
			t.Log(err)
		}

		var doc map[string][]struct{}
		if err := json.Unmarshal(bodyBuf, &doc); err != nil {
			t.Error("it should return an object")
			t.Log(err)
		}

		if val, hasKey := doc["data"]; !hasKey {
			t.Error("it should return a document with member called data")
		} else if val == nil {
			t.Error(`it should return a non nil (null) value for top level member data ie: {"data": []} not {"data": null}`)
		}
	})

	t.Run("When GET resource collection", func(t *testing.T) {
		// Setup
		req, err := http.NewRequest(http.MethodGet, "/resource", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Accept", jsonapi.ContentType)
		req.Header.Set("Content-Type", jsonapi.ContentType)

		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux

		var calledHandler bool
		mux.HandleFetchCollection("resource", func(res jsonapi.FetchCollectionResponder, req *http.Request) {
			calledHandler = true

			type NamedThing struct {
				Name string `json:"name"`
			}
			res.AppendData("0", NamedThing{"foo"}, nil, nil, nil)
			res.AppendData("1", NamedThing{"bar"}, nil, nil, nil)
			res.AppendData("2", NamedThing{"baz"}, nil, nil, nil)
		})
		// Run
		mux.ServeHTTP(res, req)

		// Test Expectaions
		if !calledHandler {
			t.Error(`it should call handler`)
		}

		result := res.Result()

		bodyBuf, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Error("it should return valid json")
			t.Log(err)
		}

		var doc map[string]interface{}
		if err := json.Unmarshal(bodyBuf, &doc); err != nil {
			t.Error("it should return an object")
			t.Log(err)
		}

		dataMemberValue, hasDataMember := doc["data"]
		if !hasDataMember {
			t.Error("it should return a document with member called data")
		}
		dataMemberValueArray, dataMemberValueIsObject := dataMemberValue.([]interface{})
		if !dataMemberValueIsObject {
			t.Error(`it should have a document member called data of type array`)
		}
		if len(dataMemberValueArray) != 3 {
			t.Errorf("it should return three values")
			t.Log(dataMemberValueArray)
		}

		for i, arrayElement := range dataMemberValueArray {
			resourceObject, isObject := arrayElement.(map[string]interface{})
			if !isObject {
				t.Error("it's data member should not have any elements that are not objects")
				t.Logf("element %d is not an object", i)
				continue
			}

			if _, hasID := resourceObject["id"]; !hasID {
				t.Error("it's data member should not have any elements without a member called id")
				t.Logf("element %d is missing member 'id'", i)
			}

			typeName, hasType := resourceObject["type"]
			if !hasType {
				t.Error("it's data member should not have any elements without a member called type")
				t.Logf("element %d is missing member 'type'", i)
			}

			if typeNameVal, isString := typeName.(string); !isString || typeNameVal != "resource" {
				t.Error(`it's elements member should have type value equal to "resource"`)
				t.Logf("element %d member has member 'type' with go type and value: %T %q", i, typeNameVal, typeNameVal)
			}

			attributes, attributesID := resourceObject["attributes"]
			if !attributesID {
				t.Error("it's data member should not have any elements without a member called id")
				t.Logf("element %d is missing member 'id'", i)
			}

			if _, isObject := attributes.(map[string]interface{}); !isObject {
				t.Error("it's data[].attribute should be an object")
			}
		}
	})

	t.Run("When GET unknown resource", func(t *testing.T) {
		// Setup
		req, err := http.NewRequest(http.MethodGet, "/unknown", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Accept", jsonapi.ContentType)
		req.Header.Set("Content-Type", jsonapi.ContentType)

		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux

		mux.HandleFetchCollection("resource", func(res jsonapi.FetchCollectionResponder, req *http.Request) {})

		// Run
		mux.ServeHTTP(res, req)

		// Test Expectaions
		result := res.Result()

		if result.StatusCode != http.StatusNotFound {
			t.Error("it should have not found status")
		}
	})

	t.Run("When GET single resource", func(t *testing.T) {
		// Setup
		req, err := http.NewRequest(http.MethodGet, "/resource/n", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Accept", jsonapi.ContentType)
		req.Header.Set("Content-Type", jsonapi.ContentType)

		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux

		var passedIDStr string
		mux.HandleFetchOne("resource", func(res jsonapi.FetchOneResonder, req *http.Request, idStr string) {
			passedIDStr = idStr

			type NamedThing struct {
				Name string `json:"name"`
			}
			res.SetData("n", NamedThing{"foo"}, nil, nil, nil)
		})

		// Run
		mux.ServeHTTP(res, req)

		// Test Expectaions
		// result := res.Result()

		if passedIDStr != "n" {
			t.Error(`it should have passed idStr "n"`)
			t.Log(passedIDStr)
		}
	})
}

func Test_ServeHTTP_AppendError(t *testing.T) {
	doc := jsonapi.TopLevelDocument{}

	doc.AppendError(nil)
	if len(doc.Errors) != 0 {
		t.Error("should not append nil error")
	}

	doc.AppendError(errors.New("lemon"))
	if len(doc.Errors) != 1 {
		t.Error("should append non nil error")
	}
}
