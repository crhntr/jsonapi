package jsonapi

import (
	"encoding/json"
	"fmt"
)

// Relationship represents a single named relationship either a to many or to
// one relationship and it may contain links and meta.
type Relationship struct {
	Data ResourceLinkage `json:"data,omitempty"`

	Links Links `json:"link,omitempty"`
	Meta  Meta  `json:"meta,omitempty"`
}

// Relationships represents a “relationships object” members of this object
// represent refreences from the resource object to other resource object.
type Relationships map[string]Relationship

// SetToOne safely sets a ToOne foreign {id,type] for a related resource.
func (rels Relationships) SetToOne(relationshipName, resourceType, id string, meta Meta) error {
	rel := rels[relationshipName]
	rel.Data.ToOne = Identifier{id, resourceType}
	rels[relationshipName] = rel

	if rel.Data.ToMany != nil {
		rel.Data.ToMany = nil
		return fmt.Errorf("to many relationship already set for %q", relationshipName)
	}
	return nil
}

// AppendToMany safely sets a ToMany foreign {id,type} for a related resource.
func (rels Relationships) AppendToMany(relationshipName, resourceType, id string, meta Meta) error {
	rel := rels[relationshipName]
	rel.Data.ToMany = append(rel.Data.ToMany, Identifier{id, resourceType})
	rels[relationshipName] = rel

	if rel.Data.ToOne.ID != "" || rel.Data.ToOne.Type != "" {
		rel.Data.ToOne.ID, rel.Data.ToOne.Type = "", ""
		return fmt.Errorf("to one relationship already set for %q", relationshipName)
	}
	return nil
}

// ResourceLinkage handles the duality, to-one or to-many, of a relationship
// object data member.
type ResourceLinkage struct {
	ToOne  Identifier
	ToMany []Identifier
}

// IsToMany checks if a ToMany Resource Identifier has been set.
func (linkage ResourceLinkage) IsToMany() bool {
	return linkage.ToMany != nil
}

// MarshalJSON handles proper encoding of json data representing a
// ResourceLinkage. It preferes two many relationships
func (linkage ResourceLinkage) MarshalJSON() ([]byte, error) {
	if linkage.IsToMany() {
		return json.Marshal(linkage.ToMany)
	}
	if linkage.ToOne.ID == "" && linkage.ToOne.Type == "" {
		return json.Marshal(nil)
	}
	return json.Marshal(linkage.ToOne)
}

// UnmarshalJSON handles proper decoding of json data representing a
// ResourceLinkage.
func (linkage *ResourceLinkage) UnmarshalJSON(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	if buf[0] == '[' {
		return json.Unmarshal(buf, &linkage.ToMany)
	}
	return json.Unmarshal(buf, &linkage.ToOne)
}
