package jsonapi

import "encoding/json"

type Relationship struct {
	Data ResourceLinkage `json"data,omitempty"`

	// TODO: add Links
	// TODO: add Meta
}

type Relationships map[string]Relationship

func (rels Relationships) SetToOne(relationshipName, resourceType, id string, meta Meta) error {
	rel := rels[relationshipName]
	rel.Data.ToMany = nil
	rel.Data.ToOne = Identifier{id, resourceType}
	rels[relationshipName] = rel
	return nil
}
func (rels Relationships) AppendToMany(relationshipName, resourceType, id string, meta Meta) error {
	rel := rels[relationshipName]
	rel.Data.ToOne.ID, rel.Data.ToOne.Type = "", ""
	rel.Data.ToMany = append(rel.Data.ToMany, Identifier{id, resourceType})
	rels[relationshipName] = rel
	return nil
}

type ResourceLinkage struct {
	ToOne  Identifier
	ToMany []Identifier
}

func (linkage ResourceLinkage) IsToMany() bool {
	return linkage.ToMany != nil
}

func (linkage ResourceLinkage) MarshalJSON() ([]byte, error) {
	if linkage.IsToMany() {
		return json.Marshal(linkage.ToMany)
	}
	return json.Marshal(linkage.ToOne)
}

func (linkage *ResourceLinkage) UnmarshalJSON(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	if buf[0] == '[' {
		return json.Unmarshal(buf, &linkage.ToMany)
	}
	return json.Unmarshal(buf, &linkage.ToOne)
}
