package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

//Creating this middleware to send a request id to each request user sends which can help
//Support team

type ctxKey int

const (
	requestIDKey ctxKey = iota
)

const (
	requestId = "X-Request-ID"
)

func RequestId(next http.Handler) http.Handler {

	//we do http.HandlerFunc to convert a normal func to
	// be able to serve api traffic
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestId)
		if id == "" {
			id = uuid.NewString()
		}
		w.Header().Add(requestId, id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIdFromContext(ctx context.Context) string {
	requestId := ctx.Value(requestIDKey).(string)
	return requestId
}
