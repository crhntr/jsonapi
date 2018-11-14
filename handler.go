package jsonapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Logger interface {
	Log(message string)
}

type Handle struct {
	Logger Logger
}

type Meta map[string]interface{}

func (handle Handle) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	accept := req.Header.Get("Accept")
	if !strings.HasPrefix(accept, ContentType) {
		res.WriteHeader(http.StatusNotAcceptable)
		return
	}
	contentType := req.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, ContentType) {
		res.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	res.Header().Set("Content-Type", ContentType)

	type Body struct {
		Data   interface{} `json:"data,omitempty"`
		Errors interface{} `json:"errors,omitempty"`
		Meta   Meta        `json:"meta,omitempty"`
	}

	var body Body
	if err := json.NewEncoder(res).Encode(body); err != nil {
		handle.Logger.Log("encoding/json.Encode could not encode the root document and returned an error: " + err.Error())
	}
}
