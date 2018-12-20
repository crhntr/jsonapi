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

// Endpoint retrieves the endpoint from the request path added to the request
// context in the router
func Endpoint(ctx context.Context) string {
	val := ctx.Value(endpointContextKey)
	endpoint, _ := val.(string)
	return endpoint
}
