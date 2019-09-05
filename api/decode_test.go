package api_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/graymeta/gmkit/api"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	body := strings.NewReader(`{
		"name":"Piotr",
		"number": 123456
	}`)
	r, err := http.NewRequest(http.MethodPost, "/something", body)
	require.NoError(t, err)
	var obj struct {
		Name   string
		Number int
	}
	require.NoError(t, api.Decode(r, &obj))
	require.Equal(t, "Piotr", obj.Name)
	require.Equal(t, 123456, obj.Number)
}

type notvalid struct{}

func (notvalid) OK() error {
	return errors.New("not ok")
}

func TestOK(t *testing.T) {
	var obj notvalid
	r, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"something":true}`))
	require.NoError(t, err)
	err = api.Decode(r, &obj)
	require.Equal(t, "validation of decoded object failed: not ok", err.Error())
}
