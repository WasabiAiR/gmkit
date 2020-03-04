package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/graymeta/gmkit/testhelpers"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNoCache(t *testing.T) {
	r := mux.NewRouter()

	testhelpers.GET(t, NoCache(r), "/",
		testhelpers.StatusNotFound(),
		testhelpers.Resp(func(t *testing.T, w *httptest.ResponseRecorder) {
			// check the headers
			assert.NotEqual(t, "", w.Header().Get("Expires"))
			assert.Equal(t, "no-cache, private, max-age=0", w.Header().Get("Cache-Control"))
			assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
			assert.Equal(t, "0", w.Header().Get("X-Accel-Expires"))
		}),
	)
}
