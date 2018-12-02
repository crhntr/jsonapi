package jsonapi

import (
	"net/http"
)

type (
	// FetchOneFunc defines how to handle a request for a single resource.
	FetchOneFunc func(res FetchOneResonder, req *http.Request, idStr string)

	// FetchCollectionFunc defines how to handle a request for a collection of
	// resources.
	FetchCollectionFunc func(res FetchCollectionResponder, req *http.Request)

	// FetchRelatedFunc defines how to handle a request for a related resource.
	FetchRelatedFunc func(res FetchRelatedResponder, req *http.Request, id, relation string)

	// FetchRelationshipsFunc defines how to handle a request for the identities
	// of a relationship and the responder provides methods to render either a
	// to-one or to-many relationship. SetDataCollection should be called when
	// the relationship represents an empty to-many relationship.
	FetchRelationshipsFunc func(res FetchRelationshipsResponder, req *http.Request, id, relation string)

	// FetchCollectionResponder represents the 'ResponseWriter' for FetchOneFunc
	FetchCollectionResponder interface {
		DataAppender
		ErrorAppender
		Includer
	}

	// FetchOneResonder represents the 'ResponseWriter' for FetchCollectionFunc
	FetchOneResonder interface {
		DataSetter
		ErrorAppender
		Includer
	}

	// FetchRelatedResponder represents the 'ResponseWriter' for FetchRelatedFunc
	FetchRelatedResponder interface {
		DataSetter
		DataAppender
		ErrorAppender
	}

	// FetchRelationshipsResponder represents the 'ResponseWriter' for
	// FetchRelationshipsFunc
	FetchRelationshipsResponder interface {
		IdentitySetter
		IdentityAppender
		DataCollectionSetter
	}

	fetchResponder interface {
		http.ResponseWriter
		DataSetter
		DataAppender
		IdentitySetter
		IdentityAppender
		ErrorAppender
		Includer
		DataCollectionSetter
	}

	fetchHandler struct {
		one FetchOneFunc
		col FetchCollectionFunc

		related       map[string]FetchRelatedFunc
		relationships map[string]FetchRelationshipsFunc
	}
)

func (hand fetchHandler) handle(res fetchResponder, req *http.Request, _ string) {
	if req.URL.Path == "/" {
		if hand.col == nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		res.SetDataCollection()
		hand.col(res, req)
		return
	}

	var (
		id, rel string
	)
	id, req.URL.Path = shiftPath(req.URL.Path)
	if req.URL.Path == "/" {
		if hand.one == nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		hand.one(res, req, id)
		return
	}

	rel, req.URL.Path = shiftPath(req.URL.Path)

	if rel == "relationships" {
		rel, req.URL.Path = shiftPath(req.URL.Path)

		relationshipsHand, ok := hand.relationships[rel]
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		relationshipsHand(res, req, id, rel)
		return
	}

	if hand.related == nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	relatedHand, ok := hand.related[rel]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	relatedHand(res, req, id, rel)
}
