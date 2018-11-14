package jsonapi

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
)

type Logger interface {
	Log(message string)
}

type Mux struct {
	Logger Logger
}

type Meta map[string]interface{}

func (mux Mux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
		mux.Logger.Log("encoding/json.Encode could not encode the root document and returned an error: " + err.Error())
	}
}

type FetchOneResonder interface{}
type FetchManyResponder interface{}
type CreateResponder interface{}
type UpdateResponder interface{}
type DeleteResponder interface{}

type FetchOneFunc func(res FetchOneResonder, req *http.Request, idStr string)
type FetchManyFunc func(res FetchManyResponder, req *http.Request)
type CreateFunc func(res CreateResponder, req *http.Request)
type UpdateFunc func(res UpdateResponder, req *http.Request, idStr string)
type DeleteFunc func(res DeleteResponder, req *http.Request, idStr string)

func (mux *Mux) HandleFetchOne(resourceName string, fn FetchOneFunc)   {}
func (mux *Mux) HandleFetchMany(resourceName string, fn FetchManyFunc) {}
func (mux *Mux) HandleCreate(resourceName string, fn CreateFunc)       {}
func (mux *Mux) HandleUpdate(resourceName string, fn UpdateFunc)       {}
func (mux *Mux) HandleDelete(resourceName string, fn DeleteFunc)       {}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
