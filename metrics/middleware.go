package metrics

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const (
	// format for the stats
	// for example for status:
	//    api.%d.status.%d
	//    api.search.status.400

	statStatus = "api.%s.status.%d"
	statTotal  = "api.%s.total"
	statTime   = "api.%s.responsetime"
)

// Handler provides a middleware that will report how much time the wrapped handler took
func Handler(next http.Handler, endpoint string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &RecorderResponseWriter{w, 0}
		ep := strings.Replace(endpoint, ".", "_", -1)
		stop := NewTimer(fmt.Sprintf(statTime, ep))

		next.ServeHTTP(recorder, r)
		status := recorder.Status()
		stop()

		// pub the stats
		Incr(fmt.Sprintf(statStatus, ep, status), 1)
		Incr(fmt.Sprintf(statTotal, ep), 1)
	})
}

// Middleware is a metrics handler middleware. If your route contains wildcards
// (ie `/some/route/with/{id_of_object}/foo`) you should not use this middleware
// as it will create individual metrics for each value of `{id_of_object}`.
// Instead, you want to use Handler() and send it the same route that you gave
// to the router, but sanitized.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &RecorderResponseWriter{w, 0}
		ep := strings.Replace(SanitizeRoute(r.Method, r.URL.Path), ".", "_", -1)
		stop := NewTimer(fmt.Sprintf(statTime, ep))

		next.ServeHTTP(recorder, r)
		status := recorder.Status()
		stop()

		// pub the stats
		Incr(fmt.Sprintf(statStatus, ep, status), 1)
		Incr(fmt.Sprintf(statTotal, ep), 1)
	})
}

// HandlerFuncCtx is the same implementation as Handler but accepts context.Context as the first parameter
func HandlerFuncCtx(h func(context.Context, http.ResponseWriter, *http.Request), endpoint string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &RecorderResponseWriter{w, 0}
		ep := strings.Replace(endpoint, ".", "_", -1)
		stop := NewTimer(fmt.Sprintf(statTime, ep))

		h(r.Context(), recorder, r)
		status := recorder.Status()
		stop()

		// pub the stats
		Incr(fmt.Sprintf(statStatus, ep, status), 1)
		Incr(fmt.Sprintf(statTotal, ep), 1)
	})
}

// HandlerFunc is a convenience method to wrap http.HandlerFunc with the middleware
func HandlerFunc(next http.HandlerFunc, endpoint string) http.Handler {
	return Handler(next, endpoint)
}

// Recorder is a convenience interface to record the response codes
type Recorder interface {
	http.ResponseWriter
	Status() int
}

// RecorderResponseWriter is an http.ResponseWriter that keeps track of the http status code
type RecorderResponseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader keeps track of the status code
func (r *RecorderResponseWriter) WriteHeader(code int) {
	r.ResponseWriter.WriteHeader(code)
	r.status = code
}

// Flush flushes the writer
func (r *RecorderResponseWriter) Flush() {
	flusher, ok := r.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

// Status returns the status code
func (r *RecorderResponseWriter) Status() int {
	return r.status
}

// Written returns true if a status code has been seen
func (r *RecorderResponseWriter) Written() bool {
	return r.Status() != 0
}

// Hijack fulfills the Hijacker interface and allows the handler to take over the connection
func (r *RecorderResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}
