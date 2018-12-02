package jsonapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error objects provide additional information about problems encountere while proforming
// an operation.
type Error struct {
	// ID may be a unique identifier for this particular occurrence of the problem.
	ID string `json:"id,omitempty"`

	// Links may be a links object containing the following members:
	Links Links `json:"links,omitempty"`

	// About may be a link that leads to further details about this particular occurrence of the problem.
	About *Link `json:"about,omitempty"`

	// Status may be the HTTP status code applicable to this problem, expressed as a string value.
	Status int `json:"status,omitempty,string"`

	// Code may be an application-specific error code, expressed as a string value.
	Code string `json:"code,omitempty"`

	// Title may be a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem, except for purposes of localization.
	Title string `json:"title,omitempty"`

	// Detail may be a human-readable explanation specific to this occurrence of the problem. Like title, this field’s value can be localized.
	Detail string `json:"detail,omitempty"`

	// Source may be an object containing references to the source of the error, optionally including any of the following members:
	Source string `json:"source,omitempty"`

	// // Pointer may be a JSON Pointer [RFC6901] to the associated entity in the request document [e.g. "/data" for a primary data object, or "/data/attributes/title" for a specific attribute].
	Pointer string `json:"pointer,omitempty"`

	// Parameter may be a string indicating which URI query parameter caused the error.
	Parameter string `json:"parameter,omitempty"`

	// Meta may be a meta object containing non-standard meta-information about the error.
	Meta Meta `json:"meta,omitempty"`
}

// HTTPStatus returns a HTTP status code for an error
// If status error has been set, then 500 (internal server error) is returned.
func (error Error) HTTPStatus() int {
	if error.Status == 0 {
		return http.StatusInternalServerError
	}
	return error.Status
}

// Error inplements errors.Error for the jsonapi.Error type
func (error Error) Error() string {
	buf, err := json.Marshal(error)
	if err != nil {
		return fmt.Sprintf("error marshaling Error struct: %s", err.Error())
	}
	return string(buf)
}

// ErrorsPolicy is used by default to calculate response
// http status codes.
// https://jsonapi.org/format/#errors-processing
func ErrorsPolicy(errors []Error) int {
	var (
		status             = http.StatusOK
		count5XX, count4XX int
	)
	for _, err := range errors {
		if err.Status > 500 {
			count5XX++
		} else if err.Status > 400 {
			count4XX++
		}

		if err.Status > status {
			status = err.Status
		}
	}

	if count4XX > 1 {
		status = 400
	}

	if count5XX > 1 {
		status = 500
	}

	return status
}
