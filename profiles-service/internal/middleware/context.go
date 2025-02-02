package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserContextKey = contextKey("user_uid")

// Set user UID to context
func ContextSetUserUID(r *http.Request, userUID string) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, userUID)
	return r.WithContext(ctx)
}

// Get user UID from context
func ContextGetUserUID(r *http.Request) string {
	userUID, ok := r.Context().Value(UserContextKey).(string)
	if !ok {
		panic("missing user value in request context")
	}

	return userUID
}
