package icinga

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
)

func TestAPIRequestPostSuccess(t *testing.T) {
	jsonData := []byte(`
	{
		"results": [
			{
				"code": 200,
				"name": "client1.example.com!client.example.com_check_docker",
				"status": "Attributes updated.",
				"type": "Service"
			},
			{
				"code": 200,
				"name": "client.example.com!client.example.com_check_nomad_http",
				"status": "Attributes updated.",
				"type": "Service"
			}
		]
	}`)

	hler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	s := httptest.NewServer(hler)
	defer s.Close()

	icingaCfg := Config{
		BaseURL:  s.URL,
		Username: "test",
		Password: "test",
	}

	is := is.New(t)
	ic, err := icingaCfg.Client()
	is.NoErr(err)

	_, err = ic.APIRequest(http.MethodPost, "/test", nil)
	is.NoErr(err)
}

func TestAPIRequestPostFailed(t *testing.T) {
	jsonData := []byte(`
	{
		"results": [
			{
				"code": 200,
				"name": "client1.example.com!client.example.com_check_docker",
				"status": "Attributes updated.",
				"type": "Service"
			},
			{
				"code": 500,
				"name": "client.example.com!client.example.com_check_nomad_http",
				"status": "Attributes updated.",
				"type": "Service"
			}
		]
	}`)

	hler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	s := httptest.NewServer(hler)
	defer s.Close()

	icingaCfg := Config{
		BaseURL:  s.URL,
		Username: "test",
		Password: "test",
	}

	is := is.New(t)
	ic, err := icingaCfg.Client()
	is.NoErr(err)

	_, err = ic.APIRequest(http.MethodPost, "/test", nil)
	is.Err(err)
}

func TestAPIRequestPostNil(t *testing.T) {
	jsonData := []byte(`
	{
		"results": []
	}`)

	hler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	s := httptest.NewServer(hler)
	defer s.Close()

	icingaCfg := Config{
		BaseURL:  s.URL,
		Username: "test",
		Password: "test",
	}

	is := is.New(t)
	ic, err := icingaCfg.Client()
	is.NoErr(err)

	_, err = ic.APIRequest(http.MethodPost, "/test", nil)
	is.Err(err)
}

func TestAPIRequestGetNil(t *testing.T) {
	jsonData := []byte(`
	{
		"results": []
	}`)

	hler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	s := httptest.NewServer(hler)
	defer s.Close()

	icingaCfg := Config{
		BaseURL:  s.URL,
		Username: "test",
		Password: "test",
	}

	is := is.New(t)
	ic, err := icingaCfg.Client()
	is.NoErr(err)

	_, err = ic.APIRequest(http.MethodGet, "/test", nil)
	is.NoErr(err)
}

func TestAPIRequestHttpError(t *testing.T) {
	jsonData := []byte(`
{
	"results": []
}`)

	hler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	s := httptest.NewServer(hler)
	defer s.Close()

	icingaCfg := Config{
		BaseURL:  s.URL,
		Username: "test",
		Password: "test",
	}

	is := is.New(t)
	ic, err := icingaCfg.Client()
	is.NoErr(err)

	_, err = ic.APIRequest(http.MethodPost, "/test", nil)
	is.Err(err)
}

func TestAPIRequestDecodeError(t *testing.T) {
	jsonData := []byte("This should not decode")

	hler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	s := httptest.NewServer(hler)
	defer s.Close()

	icingaCfg := Config{
		BaseURL:  s.URL,
		Username: "test",
		Password: "test",
	}

	is := is.New(t)
	ic, err := icingaCfg.Client()
	is.NoErr(err)

	_, err = ic.APIRequest("GET", "/test", nil)
	is.Err(err)
}
