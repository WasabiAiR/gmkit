package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/graymeta/gmkit/api"
	"github.com/graymeta/gmkit/logger"

	"github.com/stretchr/testify/require"
)

func TestWith(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	type responseObj struct {
		Name   string `json:"name"`
		Number int    `json:"number"`
	}

	data := responseObj{
		Name:   "Piotr",
		Number: 123456,
	}
	status := http.StatusCreated

	respond := api.NewResponder(logger.Default(), "v2")
	respond.With(w, r, status, data)

	require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	require.Equal(t, "v2", w.Header().Get(api.HeaderAPIVersion))
	require.Equal(t, http.StatusCreated, w.Code)
	var decoded responseObj
	require.NoError(t, json.NewDecoder(w.Body).Decode(&decoded))
	require.Equal(t, "Piotr", decoded.Name)
	require.Equal(t, 123456, decoded.Number)
}

func TestEncodingErr(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	obj := make(chan int)

	status := http.StatusCreated

	respond := api.NewResponder(logger.Default(), "v2")
	respond.With(w, r, status, obj)

	require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	require.Equal(t, "v2", w.Header().Get(api.HeaderAPIVersion))
	require.Equal(t, http.StatusInternalServerError, w.Code)
	var decoded api.ErrorResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&decoded))
	require.Equal(t, "internal server error", decoded.Error.Message)
}

func TestRespondOK(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	respond := api.NewResponder(logger.Default(), "v2")

	respond.OK(w, r)

	require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	require.Equal(t, "v2", w.Header().Get(api.HeaderAPIVersion))
	require.Equal(t, http.StatusOK, w.Code)

	type responseObj struct {
		OK bool
	}
	var decoded responseObj
	require.NoError(t, json.NewDecoder(w.Body).Decode(&decoded))
	require.True(t, decoded.OK)
}

func TestDeprecatedMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil)
	handler.ServeHTTP(w, req)

	require.Empty(t, w.Header().Get("Warning"))

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/foo", nil)
	api.Deprecated(handler).ServeHTTP(w, req)

	require.Equal(t, `299 - "Deprecated"`, w.Header().Get("Warning"))
}
