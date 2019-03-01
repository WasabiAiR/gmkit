package icinga

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cheekybits/is"
)

func TestSetDowntimeSuccess(t *testing.T) {
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

	err = ic.SetDowntime("client1.example.com", "Host", "test-author", "Test comment", time.Now(), time.Now().Add(time.Second*30))
	is.NoErr(err)

	err = ic.SetAllDowntime("client1.example.com", "test-author", "Test comment", time.Now(), time.Now().Add(time.Second*30))
	is.NoErr(err)

	err = ic.ResetDowntime("client1.example.com", "Host")
	is.NoErr(err)

	err = ic.ResetAllDowntime("client1.example.com")
	is.NoErr(err)
}

func TestSetDowntimeFailed(t *testing.T) {
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

	err = ic.SetDowntime("client1.example.com", "Host", "test-author", "Test comment", time.Now(), time.Now().Add(time.Second*30))
	is.Err(err)

	err = ic.SetAllDowntime("client1.example.com", "test-author", "Test comment", time.Now(), time.Now().Add(time.Second*30))
	is.Err(err)

	err = ic.ResetDowntime("client1.example.com", "Host")
	is.Err(err)

	err = ic.ResetAllDowntime("client1.example.com")
	is.Err(err)
}
