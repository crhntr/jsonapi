package jsonapi

import (
	"encoding/json"
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

func Test_TopLevelDocument(t *testing.T) {
	t.Run("when an error is appended", func(t *testing.T) {
		var doc TopLevelDocument
		doc.AppendError(errors.New("some error"))
		buf, err := json.Marshal(doc)
		if err != nil {
			t.Error("it should not return an error when marshalling")
			t.Log(err)
		}

		var docMap map[string]interface{}
		if err := json.Unmarshal(buf, &docMap); err != nil {
			t.Error("it should not return an error when unmarshaling to a map")
		}

		for key := range docMap {
			if key != "errors" {
				t.Error("it should not include members other than 'errors'")
				t.Log(key)
			}
		}

		errorsValue, ok := docMap["errors"]
		if !ok {
			t.Error("the TopLevelDocument encoded as json should include an errors member")
		}

		errorsArray, ok := errorsValue.([]interface{})
		if !ok {
			t.Error("the errors array should have an errors member which is an array")
			// t.Log(reflect.TypeOf(errorsValue))
			t.Log(errorsArray)
		}
		if len(errorsArray) != 1 {
			t.Error("the errors array should a single element")
			t.Log(errorsArray)
		}
		errValue := errorsArray[0]
		errMap, ok := errValue.(map[string]interface{})
		if !ok {
			t.Error("the single element from the array should be a json object")
			t.Log(errorsArray)
		}

		for key := range errMap {
			if key != "detail" /* && key != "status" */ {
				t.Error("error object should not include members other than 'detail'")
				t.Log(key)
			}
		}

		if detailString, ok := errMap["detail"].(string); !ok || detailString != "some error" {
			t.Error("error object should have a member 'detail' with the expected error message")
			t.Log(detailString)
		}

		// if statusString, ok := errMap["status"].(string); !ok || statusString != "500" {
		// 	t.Error("error object should have a member 'status' with the expected error status")
		// 	t.Log("when no status is set, encoding suspects this means an internal server error")
		// 	t.Log(statusString)
		// }
	})
}
