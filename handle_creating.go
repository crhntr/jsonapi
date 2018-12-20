package jsonapi

import (
	"encoding/json"
	"net/http"
)

type (
	// A CreateFunc implements how a resources of a given type is to be created.
	CreateFunc func(res CreateResponder, req *http.Request)

	// CreateResponder defines what to respond to a request to create a resource.
	CreateResponder interface {
		DataSetter
		ErrorAppender
	}

	// CreateRequestData represents the request body for a creating a resource.
	CreateRequestData struct {
		Data struct {
			ID            string          `json:"id,omitempty"`
			Type          string          `json:"type"`
			Attributes    json.RawMessage `json:"attributes,omitempty"`
			Relationships Relationships   `json:"relationships,omitempty"`
		} `json:"data"`
	}
)

// // UnmarshalAttributes is a convenience method to unmarshal attributes into
// // a struct or map.
// func (data CreateRequestData) UnmarshalAttributes(attributes interface{}) error {
// 	return json.Unmarshal(data.Attributes, attributes)
// }
