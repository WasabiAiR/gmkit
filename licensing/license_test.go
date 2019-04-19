package licensing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testLicenseKeyPublic  = `AQ64EJBD52VNWG54W3PFXORDPRM4UTYV2IFQ6M6YSA3XBHUDMLSN4Y52OR5GCCYXEX5I6HMYJXGFVRVTQQTW55XSFP2GCTWXM4WGKUXHG7ROJ4RU3JXZA745R3YP4FWDSZPXCMBA65VVEUCYAHHJ7B2SOLAA====`
	testLicenseKeyPrivate = `FD7YCAYBAEFXA22DN5XHIYLJNZSXEAP7QIAACAQBANIHKYQBBIAACAKEAH7YIAAAAAFP7AYFAEBP7BQAAAAP7GP7QIAWCBB5YISCH3VK3MN3ZNW6LO5CG7CZZJHRLUQLB4Z5REBXOCPIGYXE3ZR3U5D2MEFROJP2R4OZQTOMLLDLHBBHN33PEK7UMFHNOZZMMVJOON7C4TZDJWTPSB7Z3DXQ7YLMHFS7OEYCB53LKJIFQAOOT6DVE4WAAEYQFSYTGL6DNL47YFJMZEB2KFSJET4G4GPKFLGAYIGB4XHGZZMPB7A3YLPIJV2Y6Q77N3H2T7YJY4NC2AAA====`

	testLicenseKeyPrivate2 = `FD7YCAYBAEFXA22DN5XHIYLJNZSXEAP7QIAACAQBANIHKYQBBIAACAKEAH7YIAAAAAFP7AYFAEBP7BQAAAAP7GP7QIAWCBFQAV5R437XABDVFPNAXYER3YADLG3P6GA3KHUDYRVBLUE6ZYO6FIR5BDEDWHLVJ3WKBKT4Q2HMUB372Q6ICA4FQZ72Z3OCIGJSUAWXBP75XRAYP7OXFQQEPXSTPKCSLZLBMPCPXN6WMG2K25NDGRVHUVZDAEYQEITDLAUZJXC6JIAFUZJAMGNEE4FXHXHP2LVQBOOBOLPGGAXNCTQI4M43Y4MLWLWLLKKOX2DRMB3KBYAA====`
)

//var logs = flag.Bool("logs", false, "whether or not to enable logs during testing")

func TestLicense(t *testing.T) {
	lic := LicenseMF2{
		Expiration:            time.Now().Add(1 * time.Hour),
		BaseURL:               "https://foo.example.com/",
		EnforceBaseURLCheck:   true,
		LicenseGeneratedAt:    time.Now(),
		LicenseGenerationHost: "somehost.example.com",
	}

	var licenseKey string
	t.Run("normal", func(t *testing.T) {
		var err error
		licenseKey, err = Sign(testLicenseKeyPrivate, lic)
		require.NoError(t, err)

		var fromLicense LicenseMF2
		require.NoError(t, LicenseFromKey(licenseKey, testLicenseKeyPublic, &fromLicense))

		require.True(t, fromLicense.Expiration.Equal(lic.Expiration))
		require.Equal(t, fromLicense.BaseURL, lic.BaseURL)
		require.Equal(t, fromLicense.EnforceBaseURLCheck, lic.EnforceBaseURLCheck)
		require.True(t, fromLicense.LicenseGeneratedAt.Equal(lic.LicenseGeneratedAt))
		require.Equal(t, fromLicense.LicenseGenerationHost, lic.LicenseGenerationHost)
	})

	t.Run("invalid license key", func(t *testing.T) {
		var dest LicenseMF2
		require.Error(t, LicenseFromKey("abcdefg", testLicenseKeyPublic, &dest))
	})

	t.Run("valid license key, signed by wrong key", func(t *testing.T) {
		licenseKey, err := Sign(testLicenseKeyPrivate2, lic)
		require.NoError(t, err)

		var dest LicenseMF2
		err = LicenseFromKey(licenseKey, testLicenseKeyPublic, &dest)
		require.Error(t, err)
		require.Equal(t, err, ErrLicenseInvalidSignature)
	})
}
