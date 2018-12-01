package jsonapi

import (
	"encoding/json"
	"errors"
)

// Link represents both the string and link object.
// While a String and Link Object could be represented as a string and struct
// this simplifies marshaling and unmarshalling.
type Link struct {
	String string

	Object struct {
		HREF string `json:"href"`
		Meta Meta   `json:"meta,omitempty"`
	}
}

// Empty is used to ensure that the link is set. The jsonapi does not
// document an empty link type. So both of Link's json methods, MarshalJSON and
// UnmarshalJSON, ensure that some value exists.
func (ln Link) Empty() bool {
	return ln.String == "" && ln.Object.HREF == ""
}

// MarshalJSON marshals a link as either an object or string depending on how
// it has been set. If both are set, it preferes objects.
func (ln Link) MarshalJSON() ([]byte, error) {
	if ln.Empty() {
		return nil, errors.New("a link must have a string or object value")
	}
	if ln.Object.HREF != "" {
		return json.Marshal(ln.Object)
	}
	return json.Marshal(ln.String)
}

// UnmarshalJSON unmarshals a link as either an object or string depending on
// how it has been encoded.
func (ln *Link) UnmarshalJSON(buf []byte) error {
	if len(buf) == 0 {
		return errors.New("a link must have a string or object value")
	}
	if buf[0] == '{' {
		return json.Unmarshal(buf, &ln.Object)
	}
	var name string
	err := json.Unmarshal(buf, &name)
	ln.String = name
	return err
}
