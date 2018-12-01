package jsonapi

import (
	"net/http"
)

type (
	FetchOneFunc func(res FetchOneResonder, req *http.Request, idStr string)

	FetchCollectionFunc func(res FetchCollectionResponder, req *http.Request)

	FetchRelatedFunc func(res FetchRelatedResponder, req *http.Request, id, relation string)

	FetchRelationshipsFunc func(res FetchRelationshipsResponder, req *http.Request, id, relation string)

	FetchCollectionResponder interface {
		DataAppender
		ErrorAppender
		Includer
	}

	FetchOneResonder interface {
		DataSetter
		ErrorAppender
		Includer
	}

	FetchRelatedResponder interface {
		DataSetter
		DataAppender
		ErrorAppender
	}

	FetchRelationshipsResponder interface {
		IdentifierSetter
		IdentifierAppender
	}

	fetchResponder interface {
		http.ResponseWriter
		DataSetter
		DataAppender
		IdentifierSetter
		IdentifierAppender
		ErrorAppender
		Includer
	}

	fetchHandler struct {
		one FetchOneFunc
		col FetchCollectionFunc

		related       map[string]FetchRelatedFunc
		relationships map[string]FetchRelationshipsFunc
	}
)

func (hand fetchHandler) handle(res fetchResponder, req *http.Request, resourceName string) {
	if req.URL.Path == "/" {
		if hand.col == nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

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
