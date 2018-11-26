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

		if linkage.ToMany != nil || linkage.ToOne.ID != "" || linkage.ToOne.Type != "" {
			t.Error("it should not set any identifiers")
		}
	})
	t.Run("when passed a list with to identifiers", func(t *testing.T) {
		linkage := &ResourceLinkage{}
		buf := []byte(`[{"id": "lemon", "type": "citrus"}, {"id": "orange", "type": "citrus"}]`)

		if err := linkage.UnmarshalJSON(buf); err != nil {
			t.Error("it should not return an error: " + err.Error())
		}

		if linkage.ToOne.ID != "" || linkage.ToOne.Type != "" {
			t.Error("it should set the toOne identifiers to a zero value")
		}
		if len(linkage.ToMany) != 2 {
			t.Error("it should unmarshal both identifiers")
		}
		if linkage.ToMany[0].ID != "lemon" ||
			linkage.ToMany[0].Type != "citrus" ||
			linkage.ToMany[1].ID != "orange" ||
			linkage.ToMany[1].Type != "citrus" {
			t.Error("it should properly unmarshal type and ids for ids")
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
