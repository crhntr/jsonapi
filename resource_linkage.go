package jsonapi

import "encoding/json"

type ResourceLinkage struct {
	toOne  *Identifier
	toMany []Identifier
}

func (linkage ResourceLinkage) IsToMany() bool {
	return linkage.toMany != nil
}

func (linkage ResourceLinkage) MarshalJSON() ([]byte, error) {
	if linkage.IsToMany() {
		return json.Marshal(linkage.toMany)
	}
	return json.Marshal(linkage.toOne)
}

func (linkage *ResourceLinkage) UnmarshalJSON(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	if buf[0] == '[' {
		return json.Unmarshal(buf, &linkage.toMany)
	}
	return json.Unmarshal(buf, &linkage.toOne)
}
