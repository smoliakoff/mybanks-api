package graph

import "context"

type contextKey string

const internalRequestKey contextKey = "internal_request"

func WithInternalRequest(ctx context.Context, ok bool) context.Context {
	return context.WithValue(ctx, internalRequestKey, ok)
}

func IsInternalRequest(ctx context.Context) bool {
	v, ok := ctx.Value(internalRequestKey).(bool)
	return ok && v
}
