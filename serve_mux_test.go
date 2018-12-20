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
		req, err := jsonapi.NewRequest(http.MethodGet, "/", nil)
		mustNotErr(t, err)
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
		req, err := jsonapi.NewRequest(http.MethodGet, "/", nil)
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

func TestHandle_ServeHTTP_RequestMux_Deleting(t *testing.T) {
	t.Run("When deleting", func(t *testing.T) {
		req, err := jsonapi.NewRequest(http.MethodDelete, "/resource/n", nil)
		mustNotErr(t, err)
		res := httptest.NewRecorder()

		var (
			mux              jsonapi.ServeMux
			recievedEndpoint string
		)
		mux.HandleDelete("resource", jsonapi.DeleteFunc(func(res jsonapi.DeleteResponder, req *http.Request, id string) {
			recievedEndpoint = jsonapi.Endpoint(req.Context())
		}))

		// Run
		mux.ServeHTTP(res, req)

		if recievedEndpoint != "resource" {
			t.Error("it should recieve the correct endpoint parameter")
			t.Log(recievedEndpoint)
		}
	})
}

func TestHandle_ServeHTTP_RequestMux_Fetching(t *testing.T) {
	t.Run("When fetching an empty resource collection", func(t *testing.T) {
		// Setup
		req, err := jsonapi.NewRequest(http.MethodGet, "/resource", nil)
		mustNotErr(t, err)
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

	t.Run("When fetching a resource collection", func(t *testing.T) {
		// Setup
		req, err := jsonapi.NewRequest(http.MethodGet, "/resource", nil)
		mustNotErr(t, err)
		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux

		var calledHandler bool
		mux.HandleFetchCollection("resource", func(res jsonapi.FetchCollectionResponder, req *http.Request) {
			calledHandler = true

			type NamedThing struct {
				Name string `json:"name"`
			}
			res.AppendData("resource", "0", NamedThing{"foo"}, nil, nil, nil)
			res.AppendData("resource", "1", NamedThing{"bar"}, nil, nil, nil)
			res.AppendData("resource", "2", NamedThing{"baz"}, nil, nil, nil)
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

	t.Run("When fetching an unknown resource", func(t *testing.T) {
		// Setup
		req, err := jsonapi.NewRequest(http.MethodGet, "/unknown", nil)
		mustNotErr(t, err)
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

	t.Run("When fetching a single resource", func(t *testing.T) {
		// Setup
		req, err := jsonapi.NewRequest(http.MethodGet, "/resource/n", nil)
		mustNotErr(t, err)
		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux

		var passedIDStr string
		mux.HandleFetchOne("resource", func(res jsonapi.FetchOneResonder, req *http.Request, idStr string) {
			passedIDStr = idStr

			type NamedThing struct {
				Name string `json:"name"`
			}
			res.SetData("resource", "n", NamedThing{"foo"}, nil, nil, nil)
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

	t.Run("When fetching a single resource that does not exist", func(t *testing.T) {
		// Setup
		req, err := jsonapi.NewRequest(http.MethodGet, "/resource/n", nil)
		mustNotErr(t, err)
		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux

		mux.HandleFetchOne("resource", func(res jsonapi.FetchOneResonder, req *http.Request, idStr string) {
		})

		// Run
		mux.ServeHTTP(res, req)

		result := res.Result()

		bodyBuf, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Error("it should return valid json")
			t.Log(err)
			t.Log(string(bodyBuf))
		}

		var doc map[string]interface{}
		if err := json.Unmarshal(bodyBuf, &doc); err != nil {
			t.Error("it should return an object")
			t.Log(err)
			t.Log(string(bodyBuf))
		}
		dataMemberValue, hasDataMember := doc["data"]
		if !hasDataMember {
			t.Error("it should return a document with member called data")
			t.Log(string(bodyBuf))
		}
		_, dataMemberValueIsObject := dataMemberValue.(map[string]interface{})
		if !dataMemberValueIsObject {
			t.Error(`it should have document member data that is an object`)
			t.Log(string(bodyBuf))
		}
	})
}

func TestHandle_ServeHTTP_RequestMux_Creating(t *testing.T) {
	t.Run("When creating a single resource", func(t *testing.T) {
		req, err := jsonapi.NewRequest(http.MethodPost, "/resource/n", nil)
		mustNotErr(t, err)
		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux

		var (
			recievedEndpoint string
		)
		mux.HandleCreate("resource", jsonapi.CreateFunc(func(res jsonapi.CreateResponder, req *http.Request) {
			recievedEndpoint = jsonapi.Endpoint(req.Context())
		}))

		// Run
		mux.ServeHTTP(res, req)

		if recievedEndpoint != "resource" {
			t.Error("it should recieve the correct endpoint parameter")
			t.Log(recievedEndpoint)
		}
	})

	t.Run("When a request to create resource is unsupported", func(t *testing.T) {
		req, err := jsonapi.NewRequest(http.MethodPost, "/resource", nil)
		mustNotErr(t, err)
		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux
		mux.HandleFetchOne("resource", nil) // to create resource endpoint

		// Run
		mux.ServeHTTP(res, req)

		if res.Code != http.StatusForbidden {
			t.Error("it should respond with status forbiden")
			t.Log(res.Code)
		}
	})

	t.Run("When a request with an unsupported http Method is submitted", func(t *testing.T) {
		req, err := jsonapi.NewRequest(http.MethodPut, "/resource", nil)
		mustNotErr(t, err)
		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux
		mux.HandleFetchOne("resource", nil) // to create resource endpoint

		// Run
		mux.ServeHTTP(res, req)

		if res.Code != http.StatusMethodNotAllowed {
			t.Error("it should respond with status method not allowed")
			t.Log(res.Code)
		}
	})

	t.Run("When a request returns an error", func(t *testing.T) {
		req, err := jsonapi.NewRequest(http.MethodGet, "/resource/id", nil)
		mustNotErr(t, err)
		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux
		mux.HandleFetchOne("resource", jsonapi.FetchOneFunc(func(res jsonapi.FetchOneResonder, req *http.Request, idStr string) {
			res.AppendError(jsonapi.Error{Detail: "not authorized", Status: http.StatusUnauthorized})
		}))

		// Run
		mux.ServeHTTP(res, req)

		if res.Code != http.StatusUnauthorized {
			t.Error("it should respond with status StatusUnauthorized")
			t.Log(res.Code)
			t.Log(res.Body.String())
		}
	})

	t.Run("When a PATCH request is requested", func(t *testing.T) {
		req, err := jsonapi.NewRequest(http.MethodPatch, "/resource/id", nil)
		mustNotErr(t, err)
		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux
		var callCount int
		mux.HandleUpdate("resource", jsonapi.UpdateFunc(func(res jsonapi.UpdateResponder, req *http.Request, idStr string) {
			callCount++
		}))

		// Run
		mux.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Error("it should respond with status StatusOK")
			t.Log(res.Code)
			t.Log(res.Body.String())
		}

		if callCount != 1 {
			t.Error("it should call the update handler func")
		}
	})

	t.Run("When error occurs when marshaling the top level document", func(t *testing.T) {
		req, err := jsonapi.NewRequest(http.MethodPatch, "/resource/id", nil)
		mustNotErr(t, err)
		res := httptest.NewRecorder()

		var mux jsonapi.ServeMux
		mux.HandleUpdate("resource", jsonapi.UpdateFunc(func(res jsonapi.UpdateResponder, req *http.Request, idStr string) {
			res.SetData("errAttr", "someerr", errorAttr{}, nil, nil, nil)
		}))

		// Run
		mux.ServeHTTP(res, req)

		if res.Code != http.StatusInternalServerError {
			t.Error("it should respond with status StatusOK")
			t.Log(res.Code)
			t.Log(res.Body.String())
		}
	})
}

type errorAttr struct{}

func (attr errorAttr) MarshalJSON() ([]byte, error) {
	return nil, errors.New("some-err")
}
