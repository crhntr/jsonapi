package jsonapi

import "net/http"

type (
	// DeleteFunc defines how to handle a request to delete a resource.
	DeleteFunc func(res DeleteResponder, req *http.Request, id string)

	// DeleteResponder exposes an to a DeleteFunc how it's response
	// should be like.
	DeleteResponder interface {
		ErrorAppender
	}
)
