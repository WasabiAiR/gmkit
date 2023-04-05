package icinga

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ScheduleDowntime Basic body to schedule downtime for a host
type ScheduleDowntime struct {
	Type      string `json:"type"`
	Filter    string `json:"filter"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Author    string `json:"author"`
	Comment   string `json:"comment"`
}

// SetAllDowntime will turn Downtime for the host and all the services
func (c *Client) SetAllDowntime(hostname, author, comment string, start, end time.Time) error {
	err := c.SetDowntime(hostname, "Host", author, comment, start, end)
	if err != nil {
		return fmt.Errorf("SetAllDowntime on host: %w", err)
	}
	err = c.SetDowntime(hostname, "Service", author, comment, start, end)
	if err != nil {
		return fmt.Errorf("SetAllDowntime on services: %w", err)
	}
	return nil
}

// SetDowntime will set Downtime for either a host or service
func (c *Client) SetDowntime(hostname, object, author, comment string, start, end time.Time) error {
	body := ScheduleDowntime{
		Type:      object,
		Filter:    `host.name=="` + hostname + `"`,
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		Author:    author,
		Comment:   comment,
	}

	payloadJSON, marshalErr := json.Marshal(body)
	if marshalErr != nil {
		return marshalErr
	}

	_, err := c.APIRequest(http.MethodPost, "/actions/schedule-downtime", []byte(payloadJSON))
	if err != nil {
		return err
	}

	return nil
}

// ResetDowntime Basic body to schedule downtime for a host
type ResetDowntime struct {
	Type   string `json:"type"`
	Filter string `json:"filter"`
}

// ResetAllDowntime will turn Downtime off for the host and all the services
func (c *Client) ResetAllDowntime(hostname string) error {
	err := c.ResetDowntime(hostname, "Host")
	if err != nil {
		return fmt.Errorf("ResetAllDowntime on host: %w", err)
	}
	err = c.ResetDowntime(hostname, "Service")
	if err != nil {
		return fmt.Errorf("ResetAllDowntime on services: %w", err)
	}
	return nil
}

// ResetDowntime will turn Downtime off for either a host or service
func (c *Client) ResetDowntime(hostname, object string) error {
	body := ResetDowntime{
		Type:   object,
		Filter: `host.name=="` + hostname + `"`,
	}

	payloadJSON, marshalErr := json.Marshal(body)
	if marshalErr != nil {
		return marshalErr
	}

	_, err := c.APIRequest(http.MethodPost, "/actions/remove-downtime", []byte(payloadJSON))
	if err != nil {
		return err
	}

	return nil
}
