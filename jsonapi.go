package jsonapi

import (
	"errors"
	"fmt"
)

// ContentType is used in http Headers Content-Type and Accept
const ContentType = "application/vnd.api+json"

// Identity is used in Resource Identity Objects
type Identity struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// ValidateMemberName checks if a given name is approprate for a vendor name.
// It is not used internally; however, you may want to use it in your tests.
func ValidateMemberName(name string) error {
	if len(name) == 0 {
		return errors.New("a valid member name must have at least one character")
	}

	if !globallyAllowed(rune(name[0])) {
		return errors.New("a valid member name must start with a globally allowed character")
	}

	if !globallyAllowed(rune(name[len(name)-1])) {
		return errors.New("a valid member name must end with a globally allowed character")
	}

	for _, c := range name {
		if c != '_' && c != '-' && !globallyAllowed(c) {
			return fmt.Errorf("a valid member name must only have valid characters, '%c' is not allowed", c)
		}
	}
	return nil
}

func globallyAllowed(rn rune) bool {
	return (rn >= 'a' && rn <= 'z') ||
		(rn >= 'A' && rn <= 'Z') ||
		(rn >= '0' && rn <= '9')
}
