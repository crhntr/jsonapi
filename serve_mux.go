package jsonapi

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
)

// ServeMux should be used to setup endpoints.
// it implements http.Handler. It's zero value is valid.
type ServeMux struct {
	Resources map[string]EndpointHandler
}

func (mux ServeMux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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

	if req.URL.Path == "/" {
		json.NewEncoder(res).Encode(struct{}{})
		return
	}

	var endpoint string
	endpoint, req.URL.Path = shiftPath(req.URL.Path)

	hand, found := mux.Resources[endpoint]
	if !found {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	resDoc := struct {
		http.ResponseWriter
		*TopLevelDocument
	}{ResponseWriter: res, TopLevelDocument: &TopLevelDocument{}}

	status := http.StatusOK

	switch req.Method {
	case http.MethodGet:
		hand.fetch.handle(resDoc, req, endpoint)
	case http.MethodPost:
		status = http.StatusCreated
		if hand.create == nil {
			res.WriteHeader(http.StatusForbidden)
			return
		}
		hand.create(resDoc, req, endpoint)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if len(resDoc.TopLevelDocument.Errors) != 0 {
		status = ErrorsPolicy(resDoc.TopLevelDocument.Errors)
	}

	res.WriteHeader(status)
	if err := json.NewEncoder(res).Encode(resDoc.TopLevelDocument); err != nil {
		var doc TopLevelDocument
		doc.AppendError(Error{Detail: "response could not be rendered", Status: http.StatusInternalServerError})
		json.NewEncoder(res).Encode(doc)
	}
}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func (mux *ServeMux) initResources() {
	if mux.Resources == nil {
		mux.Resources = make(map[string]EndpointHandler)
	}
}

// EndpointHandler encapsulates fetch, create, update, and delete handlers
// for a single endpoint
type EndpointHandler struct {
	// PermitClientGeneratedID bool
	fetch  fetchHandler
	create CreateFunc
}

// HandleFetchOne should be used to set and endpoint handler for
// GET `/:endpoint/:id`
func (mux *ServeMux) HandleFetchOne(endpoint string, fn FetchOneFunc) {
	mux.initResources()
	handler := mux.Resources[endpoint]
	handler.fetch.one = fn
	mux.Resources[endpoint] = handler
}

// HandleFetchCollection should be used to set and endpoint handler for
// GET `/:endpoint`
func (mux *ServeMux) HandleFetchCollection(endpoint string, fn FetchCollectionFunc) {
	mux.initResources()
	handler := mux.Resources[endpoint]
	handler.fetch.col = fn
	mux.Resources[endpoint] = handler
}

// HandleCreate should be used to set and endpoint handler for
// POST `/:endpoint`
func (mux *ServeMux) HandleCreate(endpoint string, fn CreateFunc) {
	mux.initResources()
	handler := mux.Resources[endpoint]
	handler.create = fn
	mux.Resources[endpoint] = handler
}
