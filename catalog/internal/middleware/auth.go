package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

type contextKey string

const (
	UserIDContextKey   contextKey = "UserID"
	UsernameContextKey contextKey = "Username"
)

func HeaderAuth(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Header.Get("X-User-ID")
			username := r.Header.Get("X-Username")

			if userID != "" {
				log.Info("User authenticated",
					slog.String("user_id", userID),
					slog.String("username", username))

				ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
				ctx = context.WithValue(ctx, UsernameContextKey, username)

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDContextKey).(string)
	return userID, ok
}

func GetUsername(r *http.Request) (string, bool) {
	username, ok := r.Context().Value(UsernameContextKey).(string)
	return username, ok
}
