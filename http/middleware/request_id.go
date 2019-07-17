package middleware

import (
	"context"
	"net/http"

	"github.com/graymeta/gmkit/uuid"
)

// RequestHeader is the header that marks the http request id set by the server.
const RequestHeader = "X-Http-Request-Id"

type reqID int

const reqKey reqID = iota

// RequestID sets the request id in the context to a uuid that can be traced
// throughout the servers req/resp lifecycle.
func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.TimestampUUID()
		ctx := context.WithValue(r.Context(), reqKey, id)
		w.Header().Set(RequestHeader, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// GetReqID returns any request id that had been set in the context.
func GetReqID(ctx context.Context) string {
	s, _ := ctx.Value(reqKey).(string)
	return s
}
