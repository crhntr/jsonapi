package jsonapi

import (
	"net/http"
)

type Identifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Logger interface {
	Log(message string)
}

type Meta map[string]interface{}

type Linker interface{}

// func CreateLink(req *http.Request, segments ...string) string {
// 	return fmt.Sprintf("%s://%s/%s", req.URL.Scheme, req.URL.Host, strings.Join(segments, "/"))
// }

type Linkage struct {
	// Links

	Data interface{} `json:"data"`
}

type Attributes map[string]interface{} // this should be used

type DataSetter interface {
	SetData(resourceType, id string, attributes interface{}, relationships Relationships, links Linker, meta Meta) error
}

type DataAppender interface {
	AppendData(resourceType, id string, attributes interface{}, relationships Relationships, links Linker, meta Meta) error
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

// Handler Func Types Scope Interfaces for the various endpoints these types
// allow type checking to promote conformance to {json:api} Specification
type (
	FetchOneFunc        func(res FetchOneResonder, req *http.Request, idStr string)
	FetchCollectionFunc func(res FetchCollectionResponder, req *http.Request)

	CreateFunc func(res CreateResponder, req *http.Request)
	UpdateFunc func(res UpdateResponder, req *http.Request, idStr string)
	DeleteFunc func(res DeleteResponder, req *http.Request, idStr string)

	FetchIdentifierFunc           func() // todo
	FetchIdentifierCollectionFunc func() // todo

	FetchRelationFunc           func() // todo
	FetchRelationCollectionFunc func() // todo
)

type ResourceHandler struct {
	PermitClientGeneratedID bool

	FetchOne        FetchOneFunc
	FetchCollection FetchCollectionFunc
	Create          CreateFunc
	Update          UpdateFunc
	Delete          DeleteFunc

	Relationships map[string]ResourceRelationshipHandler
}

type ResourceRelationshipHandler struct {
	FetchIdentifier           FetchIdentifierFunc
	FetchIdentifierCollection FetchIdentifierCollectionFunc

	FetchRelation           FetchRelationFunc
	FetchRelationCollection FetchRelationCollectionFunc
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

func (tld *TopLevelDocument) AppendData(resourceType, id string, attributes interface{}, relationships Relationships, links Linker, meta Meta) error {
	tld.resourceSlice = append(tld.resourceSlice, Resource{
		ID:            id,
		Type:          resourceType,
		Attributes:    attributes,
		Relationships: relationships,
	})
	return nil
}

func (tld *TopLevelDocument) SetData(resourceType, id string, attributes interface{}, relationships Relationships, links Linker, meta Meta) error {
	tld.Data = &Resource{
		ID:            id,
		Type:          resourceType,
		Attributes:    attributes,
		Relationships: relationships,
	}
	return nil
}

func (tld *TopLevelDocument) Include(resourceType, id string, attributes interface{}, links Linker, meta Meta) error {
	return nil
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
