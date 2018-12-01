package jsonapi

import (
	"bytes"
	"encoding/json"
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
	t.Run("when passed an empty struct", func(t *testing.T) {
		linkage := &ResourceLinkage{}
		buf, err := json.Marshal(linkage)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(buf, []byte(`null`)) {
			t.Error("it should marshal as 'null'")
			t.Log(string(buf))
		}
	})

	t.Run("when passed a non empty list", func(t *testing.T) {
		linkage := ResourceLinkage{ToMany: []Identifier{}}
		buf, err := json.Marshal(linkage)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(buf, []byte(`[]`)) {
			t.Error("it should return an empty list")
			t.Log(string(buf))
		}
	})

	t.Run("when passed an list", func(t *testing.T) {
		linkage := ResourceLinkage{ToMany: []Identifier{{"0", "cat"}, {"1", "cat"}, {"2", "cat"}}}
		buf, err := json.Marshal(linkage)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(buf, []byte(`[{"id":"0","type":"cat"},{"id":"1","type":"cat"},{"id":"2","type":"cat"}]`)) {
			t.Error("it should return a populated list")
			t.Log(string(buf))
		}
	})

	t.Run("when passed a single identifier", func(t *testing.T) {
		linkage := ResourceLinkage{ToOne: Identifier{"0", "cat"}}

		buf, err := json.Marshal(linkage)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(buf, []byte(`{"id":"0","type":"cat"}`)) {
			t.Error("it should return a single resource identifier object")
			t.Log(string(buf))
		}
	})
}

func TestResourceLinkage_SetToOne(t *testing.T) {
	t.Run("when relationship is empty", func(t *testing.T) {
		relationships := make(Relationships)
		if err := relationships.SetToOne("relation", "resource", "1", nil); err != nil {
			t.Error("it should not return an error")
		}
	})

	t.Run("when relationship not is empty", func(t *testing.T) {
		relationships := make(Relationships)
		relationships.AppendToMany("relation", "resource", "1", nil)

		if err := relationships.SetToOne("relation", "resource", "1", nil); err == nil {
			t.Error("it should return an error")
		}
	})
}

func TestResourceLinkage_AppendToMany(t *testing.T) {
	t.Run("when relationship is empty", func(t *testing.T) {
		relationships := make(Relationships)
		if err := relationships.AppendToMany("relation", "resource", "1", nil); err != nil {
			t.Error("it should not return an error")
		}
	})

	t.Run("when relationship not is empty", func(t *testing.T) {
		relationships := make(Relationships)
		relationships.SetToOne("relation", "resource", "1", nil)

		if err := relationships.AppendToMany("relation", "resource", "1", nil); err == nil {
			t.Error("it should return an error")
		}
	})
}
