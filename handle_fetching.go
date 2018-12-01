package jsonapi

import (
	"net/http"
)

type (
	// FetchOneFunc handles a `/:endpoint/:id` endpoint
	FetchOneFunc func(res FetchOneResonder, req *http.Request, idStr string)

	// FetchCollectionFunc handles a `/:endpoint` endpoint
	FetchCollectionFunc func(res FetchCollectionResponder, req *http.Request)

	// FetchRelatedFunc handles a `/:endpoint/:id/:relation` endpoint
	FetchRelatedFunc func(res FetchRelatedResponder, req *http.Request, id, relation string)

	// FetchRelationshipsFunc handles a `/:endpoint/:id/relationships/:relation` endpoint
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
