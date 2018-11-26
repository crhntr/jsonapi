package jsonapi

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
)

type ServeMux struct {
	Logger    Logger
	Resources map[string]ResourceHandler
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

	var resourceType string
	resourceType, req.URL.Path = shiftPath(req.URL.Path)

	resourceHandler, found := mux.Resources[resourceType]
	if !found {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	resourceHandler.callHandleFunc(res, req, resourceType)
}

func (rh ResourceHandler) callHandleFunc(res http.ResponseWriter, req *http.Request, resourceType string) {
	var document TopLevelDocument
	switch req.Method {
	case http.MethodGet:
		if req.URL.Path == "/" { // Fetch Collection
			if rh.FetchCollection == nil {
				res.WriteHeader(http.StatusMethodNotAllowed)
				// TODO: write error message
				return
			}
			rh.FetchCollection(&document, req)
			if len(document.resourceSlice) == 0 {
				document.Data = Resources{}
			} else {
				document.Data = document.resourceSlice
			}
		} else { // Fetch One
			var id string
			id, req.URL.Path = shiftPath(req.URL.Path)
			if rh.FetchOne == nil {
				res.WriteHeader(http.StatusMethodNotAllowed)
				// TODO: write error message
				return
			}
			rh.FetchOne(&document, req, id)
		}
	case http.MethodPost:
		if req.URL.Path != "/" {
			res.WriteHeader(http.StatusBadRequest)
			// TODO: write error message
			return
		}
		if rh.Create == nil {
			res.WriteHeader(http.StatusMethodNotAllowed)
			// TODO: write error message
			return
		}

		doc := struct {
			Type          string          `json:"type"`
			ID            string          `json:"id"`
			Relationships Relationships   `json:"relationships"`
			Attributes    json.RawMessage `json:"attributes"`
		}{}

		// json.NewDecoder(req.Body).DisallowUnknownFields()
		if err := json.NewDecoder(req.Body).Decode(&doc); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			// TODO: write error message
			return
		}

		if !rh.PermitClientGeneratedID && doc.ID != "" {
			res.WriteHeader(http.StatusForbidden)
			return
		}

		if doc.Type != resourceType {
			res.WriteHeader(http.StatusConflict)
			// TODO: resource handler should allow list of other acceptable types for endpoint?
			// https://jsonapi.org/format/#crud-creating-responses-409
			return
		}

		rh.Create(&document, req)
		res.WriteHeader(http.StatusCreated)
	case http.MethodPatch:
		if rh.Update == nil {
			res.WriteHeader(http.StatusMethodNotAllowed)
			// TODO: write error message
			return
		}

		var id string
		id, req.URL.Path = shiftPath(req.URL.Path)
		if id == "" {
			res.WriteHeader(http.StatusBadRequest)
			// TODO: write error message
			return
		}

		rh.Update(nil, req, id)
	case http.MethodDelete:
		if rh.Delete == nil {
			res.WriteHeader(http.StatusMethodNotAllowed)
			// TODO: write error message
			return
		}

		var id string
		id, req.URL.Path = shiftPath(req.URL.Path)
		if id == "" {
			res.WriteHeader(http.StatusBadRequest)
			// TODO: write error message
			return
		}

		rh.Delete(nil, req, id)
	}

	document.Data.setType(resourceType)
	json.NewEncoder(res).Encode(document)
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
		mux.Resources = make(map[string]ResourceHandler)
	}
}

func (mux *ServeMux) HandleFetchOne(resourceType string, fn FetchOneFunc) {
	mux.initResources()
	handler := mux.Resources[resourceType]
	handler.FetchOne = fn
	mux.Resources[resourceType] = handler
}

func (mux *ServeMux) HandleFetchCollection(resourceType string, fn FetchCollectionFunc) {
	mux.initResources()
	handler := mux.Resources[resourceType]
	handler.FetchCollection = fn
	mux.Resources[resourceType] = handler
}

func (mux *ServeMux) HandleCreate(resourceType string, fn CreateFunc) {
	mux.initResources()
	handler := mux.Resources[resourceType]
	handler.Create = fn
	mux.Resources[resourceType] = handler
}

func (mux *ServeMux) HandleUpdate(resourceType string, fn UpdateFunc) {
	mux.initResources()
	handler := mux.Resources[resourceType]
	handler.Update = fn
	mux.Resources[resourceType] = handler
}

func (mux *ServeMux) HandleDelete(resourceType string, fn DeleteFunc) {
	mux.initResources()
	handler := mux.Resources[resourceType]
	handler.Delete = fn
	mux.Resources[resourceType] = handler
}

func (mux *ServeMux) HandleRelationshipIdentifierFetch(resourceType, relationName string, fn FetchIdentifierFunc) {
}
func (mux *ServeMux) HandleRelationshipIdentifierCollectionFetch(resourceType, relationName string, fn FetchIdentifierCollectionFunc) {
}
func (mux *ServeMux) HandleRelationshipFetch(resourceType, relationName string, fn FetchOneFunc) {
}
func (mux *ServeMux) HandleRelationshipCollectionFetch(resourceType, relationName string, fn FetchCollectionFunc) {
}

func (mux *ServeMux) HandleRelationshipCreate(resourceType string, fn CreateFunc) {
	mux.initResources()
	handler := mux.Resources[resourceType]
	handler.Create = fn
	mux.Resources[resourceType] = handler
}
