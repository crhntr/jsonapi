package jsonapi

import (
	"context"
	"net/http"
)

type endpointContextKeyT int

const endpointContextKey = endpointContextKeyT(0)

func contextWithEndpointValue(req *http.Request, endpoint string) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), endpointContextKey, endpoint))
}

// Endpoint retrieves the endpoint string from the request path added to the
// context in the router. If no endpoint was set, an empty string is returned.
func Endpoint(ctx context.Context) string {
	endpoint, _ := ctx.Value(endpointContextKey).(string)
	return endpoint
}
