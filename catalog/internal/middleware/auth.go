package middleware

// import (
// 	"avito-shop-service/internal/lib/handlers/response"
// 	"context"
// 	"errors"
// 	"log/slog"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// )

// type contextKey string

// const UserIDContextKey contextKey = "UserID"

// func Auth(log *slog.Logger, secret string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			tokenString := extractToken(r)
// 			if tokenString == "" {
// 				response.RespondWithError(w, log, http.StatusUnauthorized, "Unauthorized")

// 				return
// 			}

// 			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 					return nil, errors.New("unexpected signing method")
// 				}
// 				return []byte(secret), nil
// 			})

// 			if err != nil || !token.Valid {
// 				response.RespondWithError(w, log, http.StatusUnauthorized, "Unauthorized")

// 				return
// 			}

// 			claims, ok := token.Claims.(jwt.MapClaims)
// 			if !ok {
// 				response.RespondWithError(w, log, http.StatusUnauthorized, "Unauthorized")

// 				return
// 			}

// 			exp, ok := claims["exp"].(float64)
// 			if !ok || int64(exp) < time.Now().Unix() {
// 				response.RespondWithError(w, log, http.StatusUnauthorized, "token expired")

// 				return
// 			}

// 			user, ok := claims["user"].(map[string]interface{})
// 			if !ok {
// 				response.RespondWithError(w, log, http.StatusUnauthorized, "Invalid token format")

// 				return
// 			}

// 			userID, ok := user["id"].(string)
// 			if !ok {
// 				response.RespondWithError(w, log, http.StatusUnauthorized, "Invalid user ID")

// 				http.Error(w, "Invalid user ID", http.StatusUnauthorized)
// 				return
// 			}

// 			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

// func extractToken(r *http.Request) string {
// 	authHeader := r.Header.Get("Authorization")
// 	if authHeader == "" {
// 		return ""
// 	}

// 	parts := strings.Split(authHeader, " ")
// 	if len(parts) != 2 || parts[0] != "Bearer" {
// 		return ""
// 	}

// 	return parts[1]
// }

// func GetUserID(r *http.Request) (string, bool) {
// 	userID, ok := r.Context().Value(UserIDContextKey).(string)
// 	return userID, ok
// }
