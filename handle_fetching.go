package jsonapi

import "net/http"

type (
	FetchOneFunc        func(res FetchOneResonder, req *http.Request, idStr string)
	FetchCollectionFunc func(res FetchCollectionResponder, req *http.Request)
	FetchRelatedFunc    func(res FetchRelatedResponder, req *http.Request, id, relation string)

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
		DataAppender
		ErrorAppender
		Includer
	}
)

type FetchHandler struct {
	one FetchOneFunc
	col FetchCollectionFunc

	rels map[string]FetchRelatedResponder
}

func (hand FetchHandler) handle(res fetchResponder, req *http.Request, resourceName string) {
	hand.col(res, req)
}
