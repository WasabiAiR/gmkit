package pagetoken

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"net/url"
	"strconv"
)

const (
	// TokenQueryParam is the query param that identifies the next page token.
	TokenQueryParam = "page-token"

	// LimitQueryParam is the query param that identifies with the limit.
	LimitQueryParam = "limit"

	// TotalTODO indicates that the calling code doesn't send the total number of
	// hits and should be updated in the future to do so.
	TotalTODO = -1
)

// Token is a next page token constructor.
type Token struct {
	Limit    int
	Offset   int
	limitSet bool
	Params   map[string]string
}

// NextToken generates new token from existing token and number of records
// returned from operation. If number of records is less than the limit or if the
// current offset + limit is greater than or equal tot the total records then
// there is no next page. Set totalRecords to TotalTODO to preserve backwards
// compatability until you fix your handlers to provide an accurate total.
func (t Token) NextToken(numRecords, totalRecords int) string {
	if t.Limit == 0 || numRecords < t.Limit {
		return ""
	}

	if totalRecords != TotalTODO && t.Offset+t.Limit >= totalRecords {
		return ""
	}
	newToken := t
	newToken.Offset += numRecords
	return newToken.String()
}

// GetParam returns the available value for the key provided. If none exists then "" is returned.
func (t Token) GetParam(key string) string {
	if t.Params == nil {
		return ""
	}
	return t.Params[key]
}

// String implements the fmt.Stringer interfaces to get a string representation
// of the TokenDeets.
func (t Token) String() string {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(t); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(buf.Bytes())
}

// TokenOptFn is a functional option for setting token deets with sane bounds.
type TokenOptFn func(Token) Token

// DefaultLimit applies a default limit when no value had been limitSet.
func DefaultLimit(v int) TokenOptFn {
	return func(t Token) Token {
		if t.limitSet {
			return t
		}
		t.Limit = v
		return t
	}
}

// MinLimit applies a lower bound to the token limit.
func MinLimit(min int) TokenOptFn {
	return func(t Token) Token {
		if t.Limit >= min {
			return t
		}
		t.Limit = min
		return t
	}
}

// MaxLimit applies a higher bound to the token limit.
func MaxLimit(max int) TokenOptFn {
	return func(t Token) Token {
		if t.Limit <= max {
			return t
		}
		t.Limit = max
		return t
	}
}

// Param assigns a param value to an unset param key.
func Param(key, val string) TokenOptFn {
	return func(t Token) Token {
		if t.Params == nil {
			t.Params = make(map[string]string)
		}
		if _, ok := t.Params[key]; ok {
			return t
		}
		t.Params[key] = val
		return t
	}
}

// GetTokenFromQuery provides the Token from the url.Values.
func GetTokenFromQuery(q url.Values, opts ...TokenOptFn) Token {
	pageToken := GetToken(q.Get("page-token"))
	pageToken.limitSet = true
	if pageToken.Limit == 0 {
		limit, err := strconv.ParseInt(q.Get("limit"), 10, 64)
		if err == nil {
			pageToken.Limit = int(limit)
		} else {
			pageToken.limitSet = false
		}
	}

	if pageToken.Offset == 0 {
		offset, err := strconv.ParseInt(q.Get("offset"), 10, 64)
		if err == nil {
			pageToken.Offset = int(offset)
		}
	}

	for _, o := range opts {
		pageToken = o(pageToken)
	}

	return pageToken
}

// GetToken parses a token string and returns a valid Token
func GetToken(token string) Token {
	if token == "" {
		return Token{}
	}

	b, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return Token{}
	}

	var newTok Token
	dec := gob.NewDecoder(bytes.NewReader(b))
	if err := dec.Decode(&newTok); err != nil {
		return Token{}
	}

	return newTok
}
