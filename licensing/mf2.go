package licensing

import "time"

// LicenseMF2 is a GraymMeta/Curio Platform License
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
	LicenseHost string `json:"license_host"`
	// UsageHost is the url for the usage server
	UsageHost string `json:"usage_host"`
	// PublicKey is the public key used to secure communications with the GrayMeta
	// LicenseHost and UsageHost
	PublicKey string `json:"public_key"`
	// PrivateKey is the private key used to secure communications with the GrayMeta
	// LicenseHost and UsageHost
	PrivateKey string `json:"private_key"`
	// RemoteUsageEnabled is the flag to turn on/off remote usage reporting
	RemoteUsageEnabled bool `json:"remote_usage_enabled"`
	// LicenseChecksEnabled is the flag to turn on/off remote kill/licensing checks
	LicenseChecksEnabled bool `json:"license_checks_enabled"`
}

// ExpiresAt returns the license expiration time
func (l LicenseMF2) ExpiresAt() time.Time {
	return l.Expiration
}
