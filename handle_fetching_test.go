package jsonapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestFetchHandler_handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	res := struct {
		*httptest.ResponseRecorder
		*MockDataAppender
		*MockErrorAppender
		*MockIncluder
	}{httptest.NewRecorder(), NewMockDataAppender(ctrl), NewMockErrorAppender(ctrl), NewMockIncluder(ctrl)}

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
	}
}
