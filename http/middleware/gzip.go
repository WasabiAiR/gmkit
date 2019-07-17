package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	w io.Writer
}

func (gzw gzipWriter) Write(b []byte) (int, error) {
	return gzw.w.Write(b)
}

func (gzw gzipWriter) WriterHeader(statusCode int) {
	// gzw.w.Close() (gzip writer) will write a footer even if no data has been written.
	// StatusNotModified and StatusNoContent expect an empty body so swap for generic ResponseWriter instead
	if statusCode == http.StatusNoContent || statusCode == http.StatusNotModified {
		gzw.w = gzw.ResponseWriter
	}
	gzw.ResponseWriter.WriteHeader(statusCode)
}

// Compress is middleware that implements gzip compression
func Compress(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipWriter{w: gz, ResponseWriter: w}
		h.ServeHTTP(gzw, r)
	})
}
