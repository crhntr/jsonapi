package jsonapi

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"
)

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

	switch req.Method {
	case http.MethodGet:
		hand.fetch.handle(resDoc, req, endpoint)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
	}

	if err := json.NewEncoder(res).Encode(resDoc.TopLevelDocument); err != nil {
		log.Print(err)
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

type EndpointHandler struct {
	// PermitClientGeneratedID bool

	fetch fetchHandler
}

func (mux *ServeMux) HandleFetchOne(endpoint string, fn FetchOneFunc) {
	mux.initResources()
	handler := mux.Resources[endpoint]
	handler.fetch.one = fn
	mux.Resources[endpoint] = handler
}

func (mux *ServeMux) HandleFetchCollection(endpoint string, fn FetchCollectionFunc) {
	mux.initResources()
	handler := mux.Resources[endpoint]
	handler.fetch.col = fn
	mux.Resources[endpoint] = handler
}
