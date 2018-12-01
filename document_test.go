package jsonapi

import (
	"errors"
	"testing"
)

var (
	// Test TopLevelDocument implements fetchResponder
	doc interface {
		DataSetter
		DataAppender
		IdentitySetter
		IdentityAppender
		ErrorAppender
		Includer
	} = &TopLevelDocument{}
)

func Test_ServeHTTP_AppendError(t *testing.T) {
	doc := TopLevelDocument{}

	doc.AppendError(nil)
	if len(doc.Errors) != 0 {
		t.Error("should not append nil error")
	}

	doc.AppendError(errors.New("lemon"))
	if len(doc.Errors) != 1 {
		t.Error("should append non nil error")
	}
}
