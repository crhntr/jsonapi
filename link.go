package jsonapi

import (
	"encoding/json"
	"errors"
)

type Link struct {
	String string

	Object struct {
		HREF string `json:"href"`
		Meta Meta   `json:"meta,omitempty"`
	}
}

func (ln Link) Empty() bool {
	return ln.String == "" && ln.Object.HREF == ""
}

func (ln Link) MarshalJSON() ([]byte, error) {
	if ln.Empty() {
		return nil, errors.New("a link must have a string or object value")
	}
	if ln.Object.HREF != "" {
		return json.Marshal(ln.Object)
	}
	return json.Marshal(ln.String)
}

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
