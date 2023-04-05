package icinga

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// APIResults response from icinga API
// StatusCode is the http code.  Code is the code in the body from icinga
type APIResults struct {
	StatusCode int
	Exists     bool
	Results    []struct {
		Code   float64  `json:"code"`
		Errors []string `json:"errors"`
		Status string   `json:"status"`
		Name   string   `json:"name"`
		Type   string   `json:"type"`
		Attrs  struct {
			AttrsInfo map[string]any `json:"-,"`
		} `json:"attrs"`
		Joins struct{} `json:"joins"`
		Meta  struct{} `json:"meta"`
	} `json:"results"`
}

// APIRequest can change attributes to icinga objects.  Example is active checks, downtime, notifications.
// This can also be used to create and delete new object in icinga.
func (c *Client) APIRequest(method, APICall string, jsonString []byte) (APIResults, error) {

	// Build the request
	fullURL := c.cfg.BaseURL + APICall
	request, requestErr := http.NewRequest(method, fullURL, bytes.NewBuffer(jsonString))
	if requestErr != nil {
		return APIResults{}, requestErr
	}

	if c.cfg.Username != "" && c.cfg.Password != "" {
		request.SetBasicAuth(c.cfg.Username, c.cfg.Password)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	// Execute request
	response, err := c.httpClient.Do(request)
	if err != nil {
		return APIResults{}, err
	}
	defer response.Body.Close()

	// Parse the results
	var results APIResults
	results.StatusCode = response.StatusCode
	if decodeErr := json.NewDecoder(response.Body).Decode(&results); decodeErr != nil {
		return APIResults{}, decodeErr
	}

	// On http failure return results
	if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
		return results, fmt.Errorf("http status code: %d results: %v", response.StatusCode, results.Results)
	}

	// If the results is empty and GET means object does not exists.  If other method means no changes took affect
	if len(results.Results) == 0 {
		if method == http.MethodGet {
			return results, nil
		}
		return results, errors.New("icinga API results is empty")
	}
	results.Exists = true

	// For API calls that are not a GET pull the results code from the body and look for anything failed
	if method != http.MethodGet {
		for _, result := range results.Results {
			if result.Code != 200 {
				return results, errors.New("icinga API error on one of the results")
			}
		}
	}

	return results, nil
}
