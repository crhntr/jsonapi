package jsonapi_test

import (
	"errors"
	"testing"

	"github.com/crhntr/jsonapi"
)

func Test_ServeHTTP_AppendError(t *testing.T) {
	doc := jsonapi.TopLevelDocument{}

	doc.AppendError(nil)
	if len(doc.Errors) != 0 {
		t.Error("should not append nil error")
	}

	doc.AppendError(errors.New("lemon"))
	if len(doc.Errors) != 1 {
		t.Error("should append non nil error")
	}
}
