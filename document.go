package jsonapi

type (
	// Meta can be used to include non-standard meta-information
	Meta map[string]interface{}

	// Links can be used to represent links
	Links map[string]Link
)

type typeSetter interface {
	setType(resourceType string)
}

// Resource represents a single “Resource object” and appears in a JSON:API document to
// represent a resource.
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

// Resources represents an array of “Resource objects” that appear in a JSON:API
// document to represent a collection of resources.
type Resources []Resource

func (ress Resources) setType(resourceType string) {
	for i := range ress {
		ress[i].setType(resourceType)
	}
}

// TopLevelDocument represents the standard root response for all requests.
type TopLevelDocument struct {
	Data     typeSetter `json:"data,omitempty"`
	Errors   []error    `json:"errors,omitempty"`
	Meta     Meta       `json:"meta,omitempty"`
	Included Resources  `json:"included,omitempty"`

	resourceSlice Resources
}

// SetData implements DataSetter
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

// AppendData implements DataAppender
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

// SetIdentifier implements IdentifierSetter
func (tld *TopLevelDocument) SetIdentifier(resourceType, id string) error {
	return tld.SetData(resourceType, id, nil, nil, nil, nil)
}

// AppendIdentifier implements IdentifierAppender
func (tld *TopLevelDocument) AppendIdentifier(resourceType, id string) error {
	return tld.AppendData(resourceType, id, nil, nil, nil, nil)
}

// AppendError implements ErrorAppender
func (tld *TopLevelDocument) AppendError(err error) {
	if err != nil {
		tld.Errors = append(tld.Errors, err)
	}
}

// Include implements Includer
func (tld *TopLevelDocument) Include(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error {
	tld.Included = append(tld.Included, Resource{
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
