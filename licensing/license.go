package licensing

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperboloide/lk"
)

// ErrLicenseInvalidSignature is the error returned when the license has an invalid signature
var ErrLicenseInvalidSignature = errors.New("license contains invalid signature")

// License is the minimum requirements a struct has to expose for being able to be a product license
type License interface {
	ExpiresAt() time.Time
}

// IsExpired returns true if the license is expired
func IsExpired(l License) bool {
	return time.Now().After(l.ExpiresAt())
}

// Sign signs a license and returns the base32 encoded license key
func Sign(privateKey string, l any) (string, error) {
	b, err := json.Marshal(l)
	if err != nil {
		return "", fmt.Errorf("marshaling json: %w", err)
	}

	key, err := lk.PrivateKeyFromB32String(privateKey)
	if err != nil {
		return "", fmt.Errorf("unmarshaling key: %w", err)
	}

	lic, err := lk.NewLicense(key, b)
	if err != nil {
		return "", fmt.Errorf("generating license: %w", err)
	}

	licenseB32, err := lic.ToB32String()
	if err != nil {
		return "", fmt.Errorf("transforming license to base32 string: %w", err)
	}

	return licenseB32, nil
}

// LicenseFromKey takes a licenseKey and a public key and rehydrates it into a
// the destination.
func LicenseFromKey(licenseKey, publicKey string, dest any) error {
	key, err := lk.PublicKeyFromB32String(publicKey)
	if err != nil {
		return fmt.Errorf("unpacking base32 public key: %w", err)
	}

	license, err := lk.LicenseFromB32String(licenseKey)
	if err != nil {
		return fmt.Errorf("unmarshalling license from b32 string: %w", err)
	}

	if ok, err := license.Verify(key); err != nil {
		return fmt.Errorf("verifying license: %w", err)
	} else if !ok {
		return ErrLicenseInvalidSignature
	}

	if err := json.Unmarshal(license.Data, &dest); err != nil {
		return fmt.Errorf("unmarshaling license payload json: %w", err)
	}
	return nil
}
