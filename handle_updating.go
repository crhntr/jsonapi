package jsonapi

import (
	"encoding/json"
	"net/http"
)

type (
	// UpdateFunc defines how to handle a request to create a resource.
	UpdateFunc func(res UpdateResponder, req *http.Request, id string)

	UpdateRelationshipsFunc func(res UpdateRelationshipsResponder, req *http.Request, id, relation string)

	// UpdateResponder defines what to respond to a request to create a resource.
	UpdateResponder interface {
		DataSetter
		ErrorAppender
	}

	UpdateRelationshipsResponder interface {
		IdentitySetter
		IdentityAppender

		ErrorAppender
	}

	// UpdateRequestData should be used to unmarshal update resource request
	// bodies.
	UpdateRequestData struct {
		Data struct {
			ID            string          `json:"id,omitempty"`
			Type          string          `json:"type"`
			Attributes    json.RawMessage `json:"attributes,omitempty"`
			Relationships Relationships   `json:"relationships,omitempty"`
		} `json:"data"`
	}

	updateHandler struct {
		one UpdateFunc

		relationships map[string]UpdateRelationshipsFunc
	}

	updateResponder interface {
		http.ResponseWriter

		DataSetter
		IdentitySetter
		IdentityAppender
		ErrorAppender
	}
)

// // UnmarshalAttributes is a convenience method to unmarshal attributes into
// // a struct or map.
// func (data UpdateRequestData) UnmarshalAttributes(attributes interface{}) error {
// 	return json.Unmarshal(data.Attributes, attributes)
// }

func (hand updateHandler) handle(res updateResponder, req *http.Request) {
	var (
		id, rel string
	)
	id, req.URL.Path = shiftPath(req.URL.Path)
	if req.URL.Path == "/" {
		if hand.one == nil {
			res.WriteHeader(http.StatusForbidden)
			return
		}
		hand.one(res, req, id)
		return
	}

	rel, req.URL.Path = shiftPath(req.URL.Path)

	if rel != "relationships" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	rel, req.URL.Path = shiftPath(req.URL.Path)

	relationshipsHand, ok := hand.relationships[rel]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	relationshipsHand(res, req, id, rel)
	return
}
