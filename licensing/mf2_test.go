package licensing

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPingerMF2(t *testing.T) {
	l := LicenseMF2{
		PublicKey:  "012345678901234567890123456789ab",
		PrivateKey: "ab012345678901234567890123456789",
	}

	doer := &DoerMock{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			keyID := r.Header.Get(HTTPHeaderKeyID)
			sig := r.Header.Get(HTTPHeaderSignature)
			require.Len(t, keyID, 32)
			require.NotEqual(t, "", sig)

			body, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}

			// validate the signature
			mac := hmac.New(sha256.New, []byte(l.PrivateKey))
			mac.Write([]byte(body))
			require.Equal(t, base64.StdEncoding.EncodeToString(mac.Sum(nil)), sig)

			// validate the timestamp in the body
			var reqBody PingRequestMF2
			require.NoError(t, json.Unmarshal(body, &reqBody))
			require.False(t, reqBody.CurrentTime.IsZero())
			require.InDelta(t, time.Now().Unix(), reqBody.CurrentTime.Unix(), 1)

			signedResp, err := Sign(testLicenseKeyPrivate, &PingResponseMF2{Enabled: true})
			require.NoError(t, err)
			b, err := json.Marshal(&EnvelopePingResponse{Payload: signedResp})
			require.NoError(t, err)

			return &http.Response{
				Body:       io.NopCloser(bytes.NewReader(b)),
				StatusCode: http.StatusOK,
			}, nil
		},
	}

	pinger := l.Pinger(doer, testLicenseKeyPublic)

	resp, err := pinger.Ping()
	require.NoError(t, err)
	require.True(t, resp.Enabled)
}
