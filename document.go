package jsonapi

import (
	"encoding/json"
)

type (
	// Meta can be used to include non-standard meta-information
	Meta map[string]interface{}

	// Links can be used to represent links
	Links map[string]Link
)

type typeSetter interface {
	getType(resourceType string)
}

// Resource represents a single “Resource object” and appears in a JSON:API document to
// represent a resource.
type Resource struct {
	ID   string `json:"id"`
	Type string `json:"type"`

	Attributes    interface{}   `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
}

func (res *Resource) getType(resourceType string) {
	if res.Type == "" { // a resource type may be different then the endpoint
		res.Type = resourceType
	}
}

// Resources represents an array of “Resource objects” that appear in a JSON:API
// document to represent a collection of resources.
type Resources []Resource

func (ress Resources) getType(resourceType string) {
	for i := range ress {
		ress[i].getType(resourceType)
	}
}

type topLevelMembers struct {
	Meta     Meta      `json:"meta,omitempty"`
	Included Resources `json:"included,omitempty"`
}

// TopLevelDocument represents the standard root response for all requests.
type TopLevelDocument struct {
	Data   typeSetter `json:"-"`
	Errors []Error    `json:"-"`

	resourceSlice Resources
	topLevelMembers
}

// SetDataCollection is used to ensure the top level data member is encoded
// as an empty array when it is empty
func (doc *TopLevelDocument) SetDataCollection() {
	doc.resourceSlice = make(Resources, 0)
}

// MarshalJSON marshals a link as either an object or string depending on how
// it has been set. If both are set, it preferes objects.
func (doc TopLevelDocument) MarshalJSON() ([]byte, error) {
	if len(doc.Errors) != 0 {
		return json.Marshal(struct {
			Errors []Error `json:"errors"`
			topLevelMembers
		}{doc.Errors, doc.topLevelMembers})
	}

	if doc.Data == nil {
		if doc.resourceSlice != nil {
			return json.Marshal(struct {
				Data []struct{} `json:"data"`
				topLevelMembers
			}{[]struct{}{}, doc.topLevelMembers})
		}
		return json.Marshal(struct {
			Data struct{} `json:"data"`
			topLevelMembers
		}{struct{}{}, doc.topLevelMembers})
	}

	return json.Marshal(struct {
		Data typeSetter `json:"data"`
		topLevelMembers
	}{doc.Data, doc.topLevelMembers})
}

// UnmarshalJSON unmarshals a link as either an object or string depending on
// how it has been encoded.
func (doc *TopLevelDocument) UnmarshalJSON(buf []byte) error {
	return nil
}

// SetData implements DataSetter.
func (doc *TopLevelDocument) SetData(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error {
	doc.resourceSlice = nil
	doc.Data = &Resource{
		ID:            id,
		Type:          resourceType,
		Attributes:    attributes,
		Relationships: relationships,
	}
	return nil
}

// AppendData implements DataAppender.
func (doc *TopLevelDocument) AppendData(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error {
	doc.resourceSlice = append(doc.resourceSlice, Resource{
		ID:            id,
		Type:          resourceType,
		Attributes:    attributes,
		Relationships: relationships,
	})
	doc.Data = doc.resourceSlice
	return nil
}

// SetIdentity implements IdentitySetter
func (doc *TopLevelDocument) SetIdentity(resourceType, id string) error {
	return doc.SetData(resourceType, id, nil, nil, nil, nil)
}

// AppendIdentity implements IdentityAppender
func (doc *TopLevelDocument) AppendIdentity(resourceType, id string) error {
	return doc.AppendData(resourceType, id, nil, nil, nil, nil)
}

// AppendError implements ErrorAppender appending the error as
// the detail member of an Error
func (doc *TopLevelDocument) AppendError(detail error) {
	if detail != nil {
		doc.Errors = append(doc.Errors, Error{
			Detail: detail.Error(),
		})
	}
}

// Include implements Includer.
func (doc *TopLevelDocument) Include(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error {
	doc.Included = append(doc.Included, Resource{
		ID:            id,
		Type:          resourceType,
		Attributes:    attributes,
		Relationships: relationships,
	})
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
