package icinga

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
)

func TestSetNotificationsSuccess(t *testing.T) {
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

	err = ic.SetAllNotifications("client1.example.com", true)
	is.NoErr(err)

	err = ic.SetNotifications("client1.example.com", "/objects/services", true)
	is.NoErr(err)

	err = ic.SendNotification(`host.name=="cleint1"`, "Host", "jjs", "test", true)
	is.NoErr(err)
}

func TestSetNotificationsFailed(t *testing.T) {
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

	err = ic.SetAllNotifications("client1.example.com", true)
	is.Err(err)

	err = ic.SetNotifications("client1.example.com", "/objects/services", true)
	is.Err(err)

	err = ic.SendNotification(`host.name=="cleint1"`, "Host", "jjs", "test", true)
	is.Err(err)
}

func TestSetNotificationsEmpty(t *testing.T) {
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

	err = ic.SetAllNotifications("client1.example.com", true)
	is.Err(err)

	err = ic.SetNotifications("client1.example.com", "/objects/services", true)
	is.Err(err)

	err = ic.SendNotification(`host.name=="cleint1"`, "Host", "jjs", "test", true)
	is.Err(err)
}