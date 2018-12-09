package jsonapi

import (
	"encoding/json"
	"log"
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
	case http.MethodPatch:
		hand.update.handle(resDoc, req)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if len(resDoc.TopLevelDocument.Errors) != 0 {
		status = ErrorsPolicy(resDoc.TopLevelDocument.Errors)
	}

	marshaledDoc, err := json.Marshal(resDoc.TopLevelDocument)
	if err != nil {
		status = http.StatusInternalServerError

		var doc TopLevelDocument
		doc.AppendError(Error{Detail: "response could not be rendered", Status: http.StatusInternalServerError})
		marshaledDoc, err = json.Marshal(doc)
		if err != nil {
			log.Println(`{"errors": [{"detail": "top level document response could not be encoded"}]}`, err)
		}
	}

	res.WriteHeader(status)
	res.Write(marshaledDoc)
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
	update updateHandler
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

// HandleUpdate should be used to set and endpoint handler for
// PATCH `/:endpoint/:id`
func (mux *ServeMux) HandleUpdate(endpoint string, fn UpdateFunc) {
	mux.initResources()
	handler := mux.Resources[endpoint]
	handler.update.one = fn
	mux.Resources[endpoint] = handler
}
