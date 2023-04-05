package licensing

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Canonical HTTP header names for the HMAC key ID and signature headers
const (
	HTTPHeaderKeyID     = "X-Graymeta-Key-Id"
	HTTPHeaderSignature = "X-Graymeta-Signature"
)

// LicenseMF2 is a GrayMeta/Curio Platform License
type LicenseMF2 struct {
	// Expiration is the timestamp at which the license expires
	Expiration time.Time `json:"expiration"`
	// BaseURL is the platform's base URL. We include this in the signed license to
	// prevent customers from spinning up multiple instances of the platform behind
	// a load balancer
	BaseURL string `json:"base_url"`
	// EnforceBaseURLCheck dictates whether or not we should enforce the BaseURL
	// check. This is necessary because we want to use a single license across all
	// sites hosted by GrayMeta (in which case this would be false). For our typical
	// Enterprise Terraform based customers, this will be true
	EnforceBaseURLCheck bool `json:"enforce_base_url_check"`
	// LicenseGeneratedAt is the timestamp when the license was generated
	LicenseGeneratedAt time.Time `json:"license_generated_at"`
	// LicenseGenerationHost is the hostname of the server where the license was generated
	LicenseGenerationHost string `json:"license_generation_host"`
	// LicenseHost is the host for the licensing server
	LicenseHost string `json:"license_host,omitempty"`
	// UsageHost is the url for the usage server
	UsageHost string `json:"usage_host,omitempty"`
	// PublicKey is the public key used to secure communications with the GrayMeta
	// LicenseHost and UsageHost
	PublicKey string `json:"public_key,omitempty"`
	// PrivateKey is the private key used to secure communications with the GrayMeta
	// LicenseHost and UsageHost
	PrivateKey string `json:"private_key,omitempty"`
	// RemoteUsageEnabled is the flag to turn on/off remote usage reporting
	RemoteUsageEnabled bool `json:"remote_usage_enabled"`
	// LicenseChecksEnabled is the flag to turn on/off remote kill/licensing checks
	LicenseChecksEnabled bool `json:"license_checks_enabled"`
}

// ExpiresAt returns the license expiration time
func (l LicenseMF2) ExpiresAt() time.Time {
	return l.Expiration
}

// Pinger gets a new Pinger for this license. The licensePublicKey is the public
// key used to sign the mf2 licenses themselves.
func (l LicenseMF2) Pinger(client Doer, licensePublicKey string) *PingerMF2 {
	return &PingerMF2{
		client:           client,
		license:          l,
		licensePublicKey: licensePublicKey,
	}
}

// PingResponseMF2 is the response body from the license server
type PingResponseMF2 struct {
	Enabled bool `json:"enabled"`
}

// PingRequestMF2 is the request body for a license ping
type PingRequestMF2 struct {
	CurrentTime time.Time `json:"current_time"`
}

// EnvelopePingResponse wraps the signed ping response from the server
type EnvelopePingResponse struct {
	Payload string `json:"payload"`
}

// Doer is an abstraction around a http client.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// PingerMF2 executes ping requests to the MF2 license server
type PingerMF2 struct {
	client           Doer
	license          LicenseMF2
	licensePublicKey string
}

// Ping initiates a ping request to the remote license server
func (p PingerMF2) Ping() (PingResponseMF2, error) {
	bodyBytes, err := json.Marshal(PingRequestMF2{CurrentTime: time.Now()})
	if err != nil {
		return PingResponseMF2{}, fmt.Errorf("marshaling request body: %w", err)
	}

	// namespace the ping url with the app name (mf2) in case we ever have to add
	// another product and want to reuse this single licensing server
	url := fmt.Sprintf("https://%s/mf2/ping", p.license.LicenseHost)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return PingResponseMF2{}, fmt.Errorf("constructing request: %w", err)
	}

	// sign the request, add the HTTP headers
	mac := hmac.New(sha256.New, []byte(p.license.PrivateKey))
	mac.Write([]byte(bodyBytes))
	req.Header.Set(HTTPHeaderKeyID, p.license.PublicKey)
	req.Header.Set(HTTPHeaderSignature, base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	resp, err := p.client.Do(req)
	if err != nil {
		return PingResponseMF2{}, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PingResponseMF2{}, fmt.Errorf("non-200 status received: %d", resp.StatusCode)
	}

	var envelope EnvelopePingResponse
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return PingResponseMF2{}, fmt.Errorf("decoding envelope: %w", err)
	}

	// Validate the signature of the response. The payload inside the envelope is
	// signed with the license private key (that only Graymeta has access to). This
	// guarantees that the response came from a Graymeta server and not a spoofed
	// license server.
	var pingResponse PingResponseMF2
	if err := LicenseFromKey(envelope.Payload, p.licensePublicKey, &pingResponse); err != nil {
		return PingResponseMF2{}, fmt.Errorf("extracting signed content into response: %w", err)
	}

	return pingResponse, nil
}
