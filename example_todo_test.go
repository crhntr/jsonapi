package jsonapi_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/crhntr/jsonapi"
)

func ExampleServeMux() {
	var mux jsonapi.ServeMux

	mux.HandleCreate("todo", jsonapi.CreateFunc(func(res jsonapi.CreateResponder, req *http.Request, endpoint string) {
		var body jsonapi.CreateRequestData

		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			res.AppendError(errors.New("invalid json"))
		}
	}))

	testServer := httptest.NewServer(mux)
	defer testServer.Close()



}
