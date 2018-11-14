package jsonapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Logger interface {
	Log(message string)
}

type Mux struct {
	Logger Logger

	fetchHandlers map[string]fetchMux
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

	if err := json.NewEncoder(res).Encode(struct{}{}); err != nil {
		mux.Logger.Log("encoding/json.Encode could not encode the root document and returned an error: " + err.Error())
	}
}

type ResourceIdentifier struct {
	ID           string `json:"id"`
	ResourceType string `json:"type"`
	Meta         Meta   `json:"meta,omitempty"`
}

// Linker creates a "links object" following the specification
// https://jsonapi.org/format/#document-links
type Linker interface {
	Self(url string)
	Related(url string)
	Link(name, url string)

	SelfObject(href string, meta Meta)
	RelatedObject(href string, meta Meta)
	LinkObject(name, href string, meta Meta)
}

type RelationWriter interface {
	ResponseWriter

	AppendRelation(relationshipName, resourceType, id string, mata Meta)
	SetRelation(relationshipName, resourceType, id string, mata Meta)

	SetSingular(relationshipName string)
}

type RelationshipHandlerFunc func(res RelationWriter, req *http.Request, relationship string) error

type LinksHandler func(res RelationWriter, req *http.Request, links Linker) error

type Includer interface {
	Include(resourceType, id string, meta Meta, attributes interface{}, relationshipHandler RelationshipHandlerFunc, links Linker)
}

type MetaSetter interface {
	SetMeta(key string, value interface{})
}

type ResponseWriter interface {
	Includer
	MetaSetter

	Header() http.Header

	Error(title string, status int)
}

type DataSetter interface {
	SetData(resourceType, id string, meta Meta, attributes interface{}, relationshipHandler RelationshipHandlerFunc, links Linker)
}

type DataResponder interface {
	ResponseWriter
	DataSetter
}

type DataAppender interface {
	AppendData(resourceType, id string, meta Meta, attributes interface{}, relationshipHandler RelationshipHandlerFunc, links Linker)
}

type DataSliceResponder interface {
	ResponseWriter
	DataAppender
}

type RelationResponder interface {
	MetaSetter
	DataSetter
	DataAppender

	Error(title string, status int)
}

type fetchMux struct {
	One FetchOneFunc
}

func (mux Mux) ensureMaps() {
	mux.fetchHandlers = make(map[string]fetchMux)
}

type FetchOneFunc func(res DataResponder, req *http.Request)

func (mux Mux) HandleFetchOne(resourceName string, fn FetchOneFunc) {
	mux.ensureMaps()

	mp := mux.fetchHandlers[resourceName]
	mp.One = fn
	mux.fetchHandlers[resourceName] = mp
}

type ResourceMux interface {
	FetchMany(res DataSliceResponder, req *http.Request)
	FetchRelationship(res RelationResponder, req *http.Request, relationshipName string)
}
