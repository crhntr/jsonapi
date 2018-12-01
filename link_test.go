package jsonapi

import "testing"

func TestLink_Empty(t *testing.T) {
	t.Run("when link has the zero value", func(t *testing.T) {
		var link Link
		if link.Empty() != true {
			t.Error("it should return true")
		}
	})

	t.Run("when link has a string link", func(t *testing.T) {
		link := Link{String: "https://example.com/resource"}
		if link.Empty() == true {
			t.Error("it should return false")
		}
	})

	t.Run("when link has an object link", func(t *testing.T) {
		var link Link
		link.Object.HREF = "https://example.com/resource"
		link.Object.Meta = map[string]interface{}{"test": "value"}

		if link.Empty() == true {
			t.Error("it should return false")
		}
	})
}

func TestLink_MarshalJSON(t *testing.T) {
	t.Run("when link has the zero value", func(t *testing.T) {
		var link Link
		if _, err := link.MarshalJSON(); err == nil {
			t.Error("it should return an error")
		}
	})

	t.Run("when link has a string value", func(t *testing.T) {
		var link Link
		link.String = "https://example.com/resource"

		if buf, err := link.MarshalJSON(); err != nil {
			t.Error("it should not return an error")
		} else if string(buf) != `"https://example.com/resource"` {
			t.Error("it should marshal a string")
			t.Log(string(buf))
		}
	})

	t.Run("when link has an object value", func(t *testing.T) {
		var link Link
		link.Object.HREF = "https://example.com/resource"

		if buf, err := link.MarshalJSON(); err != nil {
			t.Error("it should not return an error")
		} else if string(buf) != `{"href":"https://example.com/resource"}` {
			t.Error("it should marshal an object")
			t.Log(string(buf))
		}
	})
}

func TestLink_UnmarshalJSON(t *testing.T) {
	t.Run("when link has the zero value", func(t *testing.T) {
		var link Link
		if err := link.UnmarshalJSON(nil); err == nil {
			t.Error("it should not return an error")
		}
	})

	t.Run("when link has as string value", func(t *testing.T) {
		var link Link
		if err := link.UnmarshalJSON([]byte(`"https://example.com/resource"`)); err != nil {
			t.Error("it should return an error")
		} else if link.String != "https://example.com/resource" {
			t.Error("it should set the string attribute")
			t.Log(link.String)
		}
	})

	t.Run("when link has as string value", func(t *testing.T) {
		var link Link
		if err := link.UnmarshalJSON([]byte(`{"href":"https://example.com/resource"}`)); err != nil {
			t.Error("it should return an error")
		} else if link.Object.HREF != "https://example.com/resource" {
			t.Error("it should set the object attribute")
			t.Log(link.Object.HREF)
		}
	})
}
