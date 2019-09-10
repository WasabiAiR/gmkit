package httpc

import "net/http"

// AuthFn adds authorization to an http request.
type AuthFn func(*http.Request) *http.Request

// BasicAuth sets the basic authFn on the request.
func BasicAuth(user, pass string) AuthFn {
	return func(r *http.Request) *http.Request {
		r.SetBasicAuth(user, pass)
		return r
	}
}

// HeaderAuth sets the header credential on each request.
func HeaderAuth(header, credential string) AuthFn {
	return func(r *http.Request) *http.Request {
		r.Header.Add(header, credential)
		return r
	}
}

// BearerTokenAuth sets the token authentication on the request.
func BearerTokenAuth(token string) AuthFn {
	return HeaderAuth("Authorization", "Bearer "+token)
}
