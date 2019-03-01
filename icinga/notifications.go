package icinga

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Notifications Basic body to disable notifications for a host or service
type Notifications struct {
	Filter string `json:"filter"`
	Attrs  struct {
		EnableNotificationss bool `json:"enable_notifications"`
	} `json:"attrs"`
}

// SetAllNotifications will turn notifications on or off for the host and all the services
func (c *Client) SetAllNotifications(hostname string, check bool) error {
	err := c.SetNotifications(hostname, "/objects/hosts", check)
	if err != nil {
		return errors.Wrap(err, "SetNotifications on host")
	}
	err = c.SetNotifications(hostname, "/objects/services", check)
	if err != nil {
		return errors.Wrap(err, "SetNotifications on services")
	}
	return nil
}

// SetNotifications will set notifications on or off for either a host or service
func (c *Client) SetNotifications(hostname, path string, check bool) error {
	var body Notifications
	body.Filter = `host.name=="` + hostname + `"`
	body.Attrs.EnableNotificationss = check

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

// CustomNotification body to send custom notifications for a host or service
type CustomNotification struct {
	Filter  string `json:"filter"`
	Type    string `json:"type"`
	Author  string `json:"author"`
	Comment string `json:"comment"`
	Force   bool   `json:"force"`
}

// SendNotification will send a custom notification out
func (c *Client) SendNotification(filter, checkType, author, message string, force bool) error {
	body := CustomNotification{
		Filter:  filter,
		Type:    checkType,
		Author:  author,
		Comment: message,
		Force:   force,
	}

	payloadJSON, marshalErr := json.Marshal(body)
	if marshalErr != nil {
		return marshalErr
	}

	_, err := c.APIRequest("POST", "/actions/send-custom-notification", []byte(payloadJSON))
	if err != nil {
		return err
	}

	return nil
}
