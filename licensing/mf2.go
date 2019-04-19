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
}

// ExpiresAt returns the license expiration time
func (l LicenseMF2) ExpiresAt() time.Time {
	return l.Expiration
}
