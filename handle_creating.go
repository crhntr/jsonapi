package jsonapi

import (
	"encoding/json"
	"net/http"
)

type (
	// CreateFunc defines how to handle a request to create a resource.
	CreateFunc func(res CreateResponder, req *http.Request)

	// CreateResponder defines what to respond to a request to create a resource.
	CreateResponder interface {
		DataSetter
		ErrorAppender
	}

	// CreateRequestData should be used to unmarshal create resource request
	// bodies.
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
