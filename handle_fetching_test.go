package jsonapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestFetchHandler_handle(t *testing.T) {
	type Response struct {
		*httptest.ResponseRecorder
		*MockDataSetter
		*MockDataAppender
		*MockIdentitySetter
		*MockIdentityAppender
		*MockErrorAppender
		*MockIncluder
		*MockDataCollectionSetter
	}

	mustNotErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	mustBeCalledOnce := func(callCount int) {
		if callCount != 1 {
			t.Error(`it should call hand.relationships["rel"] once`)
			t.Log(callCount)
		}
	}

	t.Run("when collection resource is fetched", func(t *testing.T) {
		t.Run("and a handler has been set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var callCount int

			hand := fetchHandler{col: FetchCollectionFunc(func(res FetchCollectionResponder, req *http.Request) {
				callCount++
			})}

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			mustNotErr(err)
			res := Response{ResponseRecorder: httptest.NewRecorder(), MockDataCollectionSetter: NewMockDataCollectionSetter(ctrl)}

			res.MockDataCollectionSetter.EXPECT().SetDataCollection().MinTimes(1)

			hand.handle(res, req, "resource")
			mustBeCalledOnce(callCount)
		})

		t.Run("and a handler has not been set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var hand fetchHandler

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			mustNotErr(err)
			res := Response{ResponseRecorder: httptest.NewRecorder()}

			hand.handle(res, req, "resource")
			if res.Code != http.StatusNotFound {
				t.Error("it should return http status not found")
			}
		})
	})

	t.Run("when one resource is fetched", func(t *testing.T) {
		t.Run("and a handler has been set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var callCount int

			hand := fetchHandler{one: FetchOneFunc(func(res FetchOneResonder, req *http.Request, idStr string) {
				callCount++
			})}

			req, err := http.NewRequest(http.MethodGet, "/0", nil)
			mustNotErr(err)
			res := Response{ResponseRecorder: httptest.NewRecorder()}

			hand.handle(res, req, "resource")
			mustBeCalledOnce(callCount)
		})

		t.Run("and a handler has not been set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var hand fetchHandler

			req, err := http.NewRequest(http.MethodGet, "/0", nil)
			mustNotErr(err)
			res := Response{ResponseRecorder: httptest.NewRecorder()}

			hand.handle(res, req, "resource")
			if res.Code != http.StatusNotFound {
				t.Error("it should return http status not found")
			}
		})
	})

	t.Run("when related resource is fetched", func(t *testing.T) {
		t.Run("and a handler has been set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var callCount int

			hand := fetchHandler{}
			hand.related = make(map[string]FetchRelatedFunc)
			hand.related["rel"] = FetchRelatedFunc(func(res FetchRelatedResponder, req *http.Request, id, relation string) {
				callCount++
			})

			req, err := http.NewRequest(http.MethodGet, "/0/rel", nil)
			mustNotErr(err)
			res := Response{ResponseRecorder: httptest.NewRecorder()}

			hand.handle(res, req, "resource")
			mustBeCalledOnce(callCount)
		})

		t.Run("and a no resource handlers have not been set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var hand fetchHandler

			req, err := http.NewRequest(http.MethodGet, "/0/rel", nil)
			mustNotErr(err)
			res := Response{ResponseRecorder: httptest.NewRecorder()}

			hand.handle(res, req, "resource")
			if res.Code != http.StatusNotFound {
				t.Error("it should return http status not found")
			}
		})

		t.Run("and a the handler for this relation has not been set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var hand fetchHandler
			hand.related = make(map[string]FetchRelatedFunc)

			req, err := http.NewRequest(http.MethodGet, "/0/rel", nil)
			mustNotErr(err)
			res := Response{ResponseRecorder: httptest.NewRecorder()}

			hand.handle(res, req, "resource")
			if res.Code != http.StatusNotFound {
				t.Error("it should return http status not found")
			}
		})
	})

	t.Run("when relationships resource is fetched", func(t *testing.T) {
		t.Run("and a handler has been set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var callCount int

			hand := fetchHandler{}
			hand.relationships = make(map[string]FetchRelationshipsFunc)
			hand.relationships["rel"] = FetchRelationshipsFunc(func(res FetchRelationshipsResponder, req *http.Request, id, relation string) {
				callCount++
			})

			req, err := http.NewRequest(http.MethodGet, "/0/relationships/rel", nil)
			mustNotErr(err)
			res := Response{ResponseRecorder: httptest.NewRecorder()}

			hand.handle(res, req, "resource")
			mustBeCalledOnce(callCount)
		})

		t.Run("and a no resource handlers have not been set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var hand fetchHandler

			req, err := http.NewRequest(http.MethodGet, "/0/relationships/rel", nil)
			mustNotErr(err)
			res := Response{ResponseRecorder: httptest.NewRecorder()}

			hand.handle(res, req, "resource")
			if res.Code != http.StatusNotFound {
				t.Error("it should return http status not found")
			}
		})

		t.Run("and a the handler for this relationship has not been set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var hand fetchHandler
			hand.relationships = make(map[string]FetchRelationshipsFunc)

			req, err := http.NewRequest(http.MethodGet, "/0/relationships/rel", nil)
			mustNotErr(err)
			res := Response{ResponseRecorder: httptest.NewRecorder()}

			hand.handle(res, req, "resource")
			if res.Code != http.StatusNotFound {
				t.Error("it should return http status not found")
			}
		})
	})
}
