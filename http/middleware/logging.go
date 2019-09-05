package middleware

import (
	"net/url"
)

var paramsToSanitize = []string{"access_token", "code"}

// SanitizeQuery used to sanitize the query params of any secrets so they do not show up in logs.
func SanitizeQuery(vals url.Values) url.Values {
	for _, v := range paramsToSanitize {
		if vals.Get(v) != "" {
			vals.Set(v, "REDACTED")
		}
	}
	return vals
}
