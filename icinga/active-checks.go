package icinga

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ActiveCheck Basic body to disable active check on a host
type ActiveCheck struct {
	Filter string `json:"filter"`
	Attrs  struct {
		EnableActiveChecks bool `json:"enable_active_checks"`
	} `json:"attrs"`
}

// SetAllActiveChecks will turn Active check on or off for the host and all the services
func (c *Client) SetAllActiveChecks(hostname string, check bool) error {
	err := c.SetActiveChecks(hostname, "/objects/hosts", check)
	if err != nil {
		return fmt.Errorf("SetActiveCheck on host: %w", err)
	}
	err = c.SetActiveChecks(hostname, "/objects/services", check)
	if err != nil {
		return fmt.Errorf("SetActiveCheck on services: %w", err)
	}
	return nil
}

// SetActiveChecks will set Active checks on or off for either a host or service
func (c *Client) SetActiveChecks(hostname, path string, check bool) error {
	var body ActiveCheck
	body.Filter = `host.name=="` + hostname + `"`
	body.Attrs.EnableActiveChecks = check

	payloadJSON, marshalErr := json.Marshal(body)
	if marshalErr != nil {
		return marshalErr
	}

	_, err := c.APIRequest(http.MethodPost, path, []byte(payloadJSON))
	if err != nil {
		return err
	}

	return nil
}
