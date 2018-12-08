package jsonapi_test

import (
	"testing"

	"github.com/crhntr/jsonapi"
)

func TestValidateMemberName(t *testing.T) {
	t.Run("when a valid name is passed", func(t *testing.T) {
		validNames := []string{
			"resource",
			"RESOURCE",
			"resourceType",
			"resourcetype",
			"resource-type",
			"resource_type",
			"n",
			"this_is_a_resource-type",
		}

		for i, name := range validNames {
			if err := jsonapi.ValidateMemberName(name); err != nil {
				t.Errorf("it should not return an error with a valid name")
				t.Logf("index: %d", i)
				t.Logf("name: %s", name)
				t.Logf("err: %s", err)
			}
		}
	})

	t.Run("when a valid name is passed", func(t *testing.T) {
		invalidNames := []string{
			"",
			" resource",
			"resource ",
			"resource type",
			"resource🙃",
			"_resource_type_",
			"-resource-type-",
			"resource-type-",
			"Resource-type-",
			"resource_type_",
			"resource+type",
			"\"resource_type\"",
		}

		for i, name := range invalidNames {
			if err := jsonapi.ValidateMemberName(name); err == nil {
				t.Errorf("it should return an error with an invalid name")
				t.Logf("index: %d", i)
				t.Logf("name: %s", name)
			}
		}
	})
}

func mustNotErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}
