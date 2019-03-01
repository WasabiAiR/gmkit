package icinga

import (
	"encoding/json"
	"net/http"
)

// Host Basic Body needed to get Host information
type Host struct {
	Filter string   `json:"filter,omitempty"`
	Attrs  []string `json:"attrs,omitempty"`
}

// HostExist will check if the host exists
func (c *Client) HostExist(hostname string) (bool, error) {
	var body Host
	body.Filter = `host.name=="` + hostname + `"`

	payloadJSON, marshalErr := json.Marshal(body)
	if marshalErr != nil {
		return false, marshalErr
	}

	result, err := c.APIRequest(http.MethodGet, "/objects/hosts", []byte(payloadJSON))
	return result.Exists, err
}
