package middleware

import (
	"context"
	"net/http"
)

const ProfilePickMiddleWarePathKey = "profilePickBasePath"

func ProfilePicMiddleware(basePath string, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ProfilePickMiddleWarePathKey, basePath)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
