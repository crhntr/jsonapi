package jsonapi

import "encoding/json"

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
