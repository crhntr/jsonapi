package jsonapi

import (
	"errors"
	"fmt"
)

const ContentType = "application/vnd.api+json"

type Identifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// ValidateMemberName expects that the name follows the stricter naming
// standards. For example, spaces are not allowed.
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
