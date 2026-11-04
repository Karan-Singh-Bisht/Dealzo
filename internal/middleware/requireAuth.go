package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/karan-singh-bisht/Dealzo-api/internal/httpx"
)

const (
	authorization = "Authorization"
)

type AuthClaims struct {
	jwt.RegisteredClaims
}

func RequireAuth(logger *slog.Logger, secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			requestId := RequestIdFromContext(ctx)
			logger.With("request_id", requestId)
			authHeader := r.Header.Get(authorization)

			if authHeader == "" {
				logger.Error("No authorization header found!")
				httpx.Error(w, http.StatusUnauthorized, "auth is required", httpx.CodeUnauthenticated)
				return
			}

			raw, ok := strings.CutPrefix(authHeader, "Bearer ")
			if !ok || raw == "" {
				logger.Error("No JWT token found!")
				httpx.Error(w, http.StatusUnauthorized, "auth is required", httpx.CodeUnauthenticated)
				return
			}

			var claims AuthClaims
			_, err := jwt.ParseWithClaims(raw, &claims, func(t *jwt.Token) (any, error) {

				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected sign in method")
				}

				return []byte(secret), nil
			}, jwt.WithValidMethods([]string{"HS256"}))

			if err != nil {
				logger.Info("token rejected", "err", err)
				httpx.Error(w, http.StatusUnauthorized, "token is invalid or expired", httpx.CodeUnauthenticated)
				return
			}

			ctxWithUserID := context.WithValue(ctx, userIDKey, claims.Subject)

			next.ServeHTTP(w, r.WithContext(ctxWithUserID))
		})
	}
}
