package jsonapi

import "testing"

func TestValidateMemberName(t *testing.T) {
	t.Run("when a valid name is passed", func(t *testing.T) {
		validNames := []string{
			"resource",
			"RESOURCE",
			"resourceName",
			"ResourceName",
			"resource-name",
			"resource_name",
			"n",
			"this_is_a_resource-name",
		}

		for i, name := range validNames {
			if err := ValidateMemberName(name); err != nil {
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
			"resource name",
			"resource🙃",
			"_resource_name_",
			"-resource-name-",
			"resource-name-",
			"Resource-name-",
			"resource_name_",
			"resource+name",
			"\"resource_name\"",
		}

		for i, name := range invalidNames {
			if err := ValidateMemberName(name); err == nil {
				t.Errorf("it should return an error with an invalid name")
				t.Logf("index: %d", i)
				t.Logf("name: %s", name)
			}
		}
	})
}
