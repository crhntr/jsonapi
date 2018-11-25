package jsonapi

import (
	"testing"
)

func TestResourceLinkage_UnmarshalJSON(t *testing.T) {
	t.Run("when passed an empty byte slice", func(t *testing.T) {
		linkage := &ResourceLinkage{}

		if err := linkage.UnmarshalJSON(nil); err != nil {
			t.Error("it should not return an error: " + err.Error())
		}
	})
	t.Run("when passed an list", func(t *testing.T) {
		linkage := &ResourceLinkage{}
		buf := []byte(`[{"id": "lemon", "type": "citrus"}, {"id": "orange", "type": "citrus"}]`)

		if err := linkage.UnmarshalJSON(buf); err != nil {
			t.Error("it should not return an error: " + err.Error())
		}
	})
	t.Run("when passed a single identifier", func(t *testing.T) {
		linkage := &ResourceLinkage{}
		buf := []byte(`{"id": "lemon", "type": "citrus"}`)

		if err := linkage.UnmarshalJSON(buf); err != nil {
			t.Error("it should not return an error: " + err.Error())
		}
	})
}

func TestResourceLinkage_MarshalJSON(t *testing.T) {
	t.Run("when passed an empty struct", func(t *testing.T) {})
	t.Run("when passed an list", func(t *testing.T) {})
	t.Run("when passed a single identifier", func(t *testing.T) {})
}
