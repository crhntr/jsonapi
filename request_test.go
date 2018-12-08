package jsonapi

import (
	"testing"
)

func TestNewRequest(t *testing.T) {
	_, err := NewRequest("bad method", "/", nil)
	if err == nil {
		t.Error("should fail")
	}
}
