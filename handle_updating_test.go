package jsonapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUpdateHandler_handle(t *testing.T) {
	type Response struct {
		*httptest.ResponseRecorder
		*MockDataSetter
		*MockIdentitySetter
		*MockIdentityAppender
		*MockErrorAppender
	}

	mustNotErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("when a resource that exists is updated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var callCount int
		hand := updateHandler{one: UpdateFunc(func(res UpdateResponder, req *http.Request, id string) {
			callCount++
		})}

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		mustNotErr(err)
		res := Response{ResponseRecorder: httptest.NewRecorder()}

		hand.handle(res, req)
		if res.Code != http.StatusOK {
			t.Error("it should respond with ok status code")
		}
		if callCount == 0 {
			t.Error("it should call the handler")
		}
	})

	t.Run("when a resource that does not exist is updated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hand := updateHandler{}

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		mustNotErr(err)
		res := Response{ResponseRecorder: httptest.NewRecorder()}

		hand.handle(res, req)
		if res.Code != http.StatusForbidden {
			t.Error("it should respond with forbidden status code")
		}
	})

	t.Run("when a updating a handled relationship", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var callCount int
		hand := updateHandler{relationships: map[string]UpdateRelationshipsFunc{"other": func(res UpdateRelationshipsResponder, req *http.Request, id, relation string) {
			callCount++
		}}}

		req, err := http.NewRequest(http.MethodGet, "/some-id/relationships/other", nil)
		mustNotErr(err)
		res := Response{ResponseRecorder: httptest.NewRecorder()}

		hand.handle(res, req)
		if res.Code != http.StatusOK {
			t.Error("it should respond with OK status code")
		}
		if callCount == 0 {
			t.Error("it should call the handler")
		}
	})

	t.Run("when a updating a not handled relationship", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hand := updateHandler{}

		req, err := http.NewRequest(http.MethodGet, "/some-id/relationships/other", nil)
		mustNotErr(err)
		res := Response{ResponseRecorder: httptest.NewRecorder()}

		hand.handle(res, req)
		if res.Code != http.StatusNotFound {
			t.Error("it should respond with not found status code")
		}
	})

	t.Run("when a update relationship path does not have relationships prefix", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hand := updateHandler{}

		req, err := http.NewRequest(http.MethodGet, "/some-id/something-else/other", nil)
		mustNotErr(err)
		res := Response{ResponseRecorder: httptest.NewRecorder()}

		hand.handle(res, req)
		if res.Code != http.StatusBadRequest {
			t.Error("it should respond with bad request status code")
		}
	})
}
