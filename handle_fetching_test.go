package jsonapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestFetchHandler_handle_FetchCollection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	res := struct {
		*httptest.ResponseRecorder
		*MockDataSetter
		*MockDataAppender
		*MockIdentifierSetter
		*MockIdentifierAppender
		*MockErrorAppender
		*MockIncluder
	}{ResponseRecorder: httptest.NewRecorder()}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	var callCount int

	hand := FetchHandler{col: FetchCollectionFunc(func(res FetchCollectionResponder, req *http.Request) {
		callCount++
	})}

	hand.handle(res, req, "resource")

	if callCount != 1 {
		t.Error("hand.col should be called once")
		t.Log(callCount)
	}
}

func TestFetchHandler_handle_FetchOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	res := struct {
		*httptest.ResponseRecorder
		*MockDataSetter
		*MockDataAppender
		*MockIdentifierSetter
		*MockIdentifierAppender
		*MockErrorAppender
		*MockIncluder
	}{ResponseRecorder: httptest.NewRecorder()}

	req, err := http.NewRequest(http.MethodGet, "/0", nil)
	if err != nil {
		t.Fatal(err)
	}

	var callCount int

	hand := FetchHandler{one: FetchOneFunc(func(res FetchOneResonder, req *http.Request, idStr string) {
		callCount++
	})}

	hand.handle(res, req, "resource")

	if callCount != 1 {
		t.Error("hand.one should be called once")
		t.Log(callCount)
	}
}

func TestFetchHandler_handle_FetchRelated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	res := struct {
		*httptest.ResponseRecorder
		*MockDataSetter
		*MockDataAppender
		*MockIdentifierSetter
		*MockIdentifierAppender
		*MockErrorAppender
		*MockIncluder
	}{ResponseRecorder: httptest.NewRecorder()}

	req, err := http.NewRequest(http.MethodGet, "/0/rel", nil)
	if err != nil {
		t.Fatal(err)
	}

	var callCount int

	hand := FetchHandler{}
	hand.related = make(map[string]FetchRelatedFunc)
	hand.related["rel"] = FetchRelatedFunc(func(res FetchRelatedResponder, req *http.Request, id, relation string) {
		callCount++
	})

	hand.handle(res, req, "resource")

	if callCount != 1 {
		t.Error(`hand.related["rel"] should be called once`)
		t.Log(callCount)
	}
}

func TestFetchHandler_handle_FetchRelationships(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	res := struct {
		*httptest.ResponseRecorder
		*MockDataSetter
		*MockDataAppender
		*MockIdentifierSetter
		*MockIdentifierAppender
		*MockErrorAppender
		*MockIncluder
	}{ResponseRecorder: httptest.NewRecorder()}

	req, err := http.NewRequest(http.MethodGet, "/0/relationships/rel", nil)
	if err != nil {
		t.Fatal(err)
	}

	var callCount int

	hand := FetchHandler{}
	hand.relationships = make(map[string]FetchRelationshipsFunc)
	hand.relationships["rel"] = FetchRelationshipsFunc(func(res FetchRelationshipsResponder, req *http.Request, id, relation string) {
		callCount++
	})

	hand.handle(res, req, "resource")

	if callCount != 1 {
		t.Error(`hand.relationships["rel"] should be called once`)
		t.Log(callCount)
	}
}
