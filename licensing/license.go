package licensing

import (
	"encoding/json"
	"time"

	"github.com/hyperboloide/lk"
	"github.com/pkg/errors"
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
func Sign(privateKey string, l interface{}) (string, error) {
	b, err := json.Marshal(l)
	if err != nil {
		return "", errors.Wrap(err, "marshaling json")
	}

	key, err := lk.PrivateKeyFromB32String(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "unmarshaling key")
	}

	lic, err := lk.NewLicense(key, b)
	if err != nil {
		return "", errors.Wrap(err, "generating license")
	}

	licenseB32, err := lic.ToB32String()
	return licenseB32, errors.Wrap(err, "transforming license to base32 string")
}

// LicenseFromKey takes a licenseKey and a public key and rehydrates it into a
// the destination.
func LicenseFromKey(licenseKey, publicKey string, dest interface{}) error {
	key, err := lk.PublicKeyFromB32String(publicKey)
	if err != nil {
		return errors.Wrap(err, "unpacking base32 public key")
	}

	license, err := lk.LicenseFromB32String(licenseKey)
	if err != nil {
		return errors.Wrap(err, "unmarshalling license from b32 string")
	}

	if ok, err := license.Verify(key); err != nil {
		return errors.Wrap(err, "verifying license")
	} else if !ok {
		return ErrLicenseInvalidSignature
	}

	if err := json.Unmarshal(license.Data, &dest); err != nil {
		return errors.Wrap(err, "unmarshaling license payload json")
	}
	return nil
}
