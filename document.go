package jsonapi

type (
	Meta  map[string]interface{}
	Links map[string]Link
)

type Attributes map[string]interface{} // this should be used

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

func (tld *TopLevelDocument) AppendData(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error {
	tld.Data = nil
	tld.resourceSlice = append(tld.resourceSlice, Resource{
		ID:            id,
		Type:          resourceType,
		Attributes:    attributes,
		Relationships: relationships,
	})
	return nil
}

func (tld *TopLevelDocument) SetData(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error {
	tld.resourceSlice = nil
	tld.Data = &Resource{
		ID:            id,
		Type:          resourceType,
		Attributes:    attributes,
		Relationships: relationships,
	}
	return nil
}

func (tld *TopLevelDocument) Include(resourceType, id string, attributes interface{}, links Links, meta Meta) error {
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
