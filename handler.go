package jsonapi

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
	"sync"
)

type Logger interface {
	Log(message string)
}

type Mux struct {
	Logger    Logger
	Resources map[string]ResourceHandler
	mut       *sync.Mutex
}

type Meta map[string]interface{}

type Linker interface{}

// func CreateLink(req *http.Request, segments ...string) string {
// 	return fmt.Sprintf("%s://%s/%s", req.URL.Scheme, req.URL.Host, strings.Join(segments, "/"))
// }

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

type Identifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Linkage struct {
	// Links

	Data interface{} `json:"data"`
}

type Attributes map[string]interface{} // this should be used

type Relationships map[string]interface{}

func (rels Relationships) SetToOne(relationshipName, resourceType, id string, meta Meta) error {
	return nil
}
func (rels Relationships) AppendToMany(relationshipName, resourceType, id string, meta Meta) error {
	return nil
}

type DataSetter interface {
	SetData(id string, attributes interface{}, relationships Relationships, links Linker, meta Meta) error
}

type DataAppender interface {
	AppendData(id string, attributes interface{}, relationships Relationships, links Linker, meta Meta) error
}

type Includer interface {
	Include(resourceType, id string, attributes interface{}, links Linker, meta Meta) error
}

type ErrorAppender interface {
	AppendError(err error)
}

type FetchOneResonder interface {
	DataSetter
	ErrorAppender
	Includer
}

type FetchCollectionResponder interface {
	DataAppender
	ErrorAppender
	Includer
}

type CreateResponder interface {
	DataSetter
}

type UpdateResponder interface{}
type DeleteResponder interface{}

type FetchOneFunc func(res FetchOneResonder, req *http.Request, idStr string)
type FetchCollectionFunc func(res FetchCollectionResponder, req *http.Request)
type CreateFunc func(res CreateResponder, req *http.Request)
type UpdateFunc func(res UpdateResponder, req *http.Request, idStr string)
type DeleteFunc func(res DeleteResponder, req *http.Request, idStr string)

func (mux *Mux) initResources() {
	if mux.mut == nil {
		mux.mut = &sync.Mutex{}
	}

	mux.mut.Lock()
	defer mux.mut.Unlock()

	if mux.Resources == nil {
		mux.Resources = make(map[string]ResourceHandler)
	}
}

func (mux *Mux) HandleFetchOne(resourceType string, fn FetchOneFunc) {
	mux.initResources()
	mux.mut.Lock()
	defer mux.mut.Unlock()
	handler := mux.Resources[resourceType]
	handler.FetchOne = fn
	mux.Resources[resourceType] = handler
}

func (mux *Mux) HandleFetchCollection(resourceType string, fn FetchCollectionFunc) {
	mux.initResources()
	mux.mut.Lock()
	defer mux.mut.Unlock()
	handler := mux.Resources[resourceType]
	handler.FetchCollection = fn
	mux.Resources[resourceType] = handler
}

func (mux *Mux) HandleCreate(resourceType string, fn CreateFunc) {
	mux.initResources()
	mux.mut.Lock()
	defer mux.mut.Unlock()
	handler := mux.Resources[resourceType]
	handler.Create = fn
	mux.Resources[resourceType] = handler
}

func (mux *Mux) HandleUpdate(resourceType string, fn UpdateFunc) {
	mux.initResources()
	mux.mut.Lock()
	defer mux.mut.Unlock()
	handler := mux.Resources[resourceType]
	handler.Update = fn
	mux.Resources[resourceType] = handler
}

func (mux *Mux) HandleDelete(resourceType string, fn DeleteFunc) {
	mux.initResources()
	mux.mut.Lock()
	defer mux.mut.Unlock()
	handler := mux.Resources[resourceType]
	handler.Delete = fn
	mux.Resources[resourceType] = handler
}

type ResourceHandler struct {
	PermitClientGeneratedID bool

	FetchOne        FetchOneFunc
	FetchCollection FetchCollectionFunc
	Create          CreateFunc
	Update          UpdateFunc
	Delete          DeleteFunc

	Relationships map[string]ResourceRelationshipHandler
}

type FetchIdentifierFunc func()           // todo
type FetchIdentifierCollectionFunc func() // todo
type FetchRelationFunc func()             // todo
type FetchRelationCollectionFunc func()   // todo

type ResourceRelationshipHandler struct {
	FetchIdentifier           FetchIdentifierFunc
	FetchIdentifierCollection FetchIdentifierCollectionFunc

	FetchRelation           FetchRelationFunc
	FetchRelationCollection FetchRelationCollectionFunc
}

func (mux *Mux) HandleFetchIdentifier(resourceType, relationName string, fn FetchIdentifierFunc) {
}
func (mux *Mux) HandleFetchIdentifierCollection(resourceType, relationName string, fn FetchIdentifierCollectionFunc) {
}
func (mux *Mux) HandleFetchRelation(resourceType, relationName string, fn FetchRelationFunc) {
}
func (mux *Mux) HandleFetchRelationCollection(resourceType, relationName string, fn FetchRelationCollectionFunc) {
}

func (mux *Mux) HandleCreate(resourceType string, fn CreateFunc) {
	mux.initResources()
	mux.mut.Lock()
	defer mux.mut.Unlock()
	handler := mux.Resources[resourceType]
	handler.Create = fn
	mux.Resources[resourceType] = handler
}

type typeSetter interface {
	setType(resourceType string)
}

type Resource struct {
	ID   string `json:"id"`
	Type string `json:"type"`

	Attributes    interface{}   `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
}

func (res *Resource) setType(resourceType string) {
	if res.Type == "" { // a resource type may be different then the endpoint
		res.Type = resourceType
	}
}

type Resources []Resource

func (ress Resources) setType(resourceType string) {
	for i := range ress {
		ress[i].setType(resourceType)
	}
}

type TopLevelDocument struct {
	Data   typeSetter `json:"data,omitempty"`
	Errors []error    `json:"errors,omitempty"`
	Meta   Meta       `json:"meta,omitempty"`

	resourceSlice Resources
}

func (tld *TopLevelDocument) AppendError(err error) {
	if err != nil {
		tld.Errors = append(tld.Errors, err)
	}
}

func (tld *TopLevelDocument) AppendData(id string, attributes interface{}, relationships Relationships, links Linker, meta Meta) error {
	tld.resourceSlice = append(tld.resourceSlice, Resource{
		ID:            id,
		Attributes:    attributes,
		Relationships: relationships,
	})
	return nil
}

func (tld *TopLevelDocument) SetData(id string, attributes interface{}, relationships Relationships, links Linker, meta Meta) error {
	tld.Data = &Resource{
		ID:            id,
		Attributes:    attributes,
		Relationships: relationships,
	}
	return nil
}

func (tld *TopLevelDocument) Include(resourceType, id string, attributes interface{}, links Linker, meta Meta) error {
	return nil
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

// func UnmarshalAttributes(req *http.Request, attributes interface{}) {
//
// }
//
// func UnmarshalToOneRelationship(req *http.Request, relationshipName string) {
//
// }
//
// func UnmarshalToManyRelationship(req *http.Request, relationshipName string) []Relationships {
//
// }
