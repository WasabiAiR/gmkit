package pagetoken

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetTokenDeets(t *testing.T) {
	t.Run("bad token returned as empty string", func(t *testing.T) {
		tests := []string{"asdfasdf", "", "anotherbadtoken", "1-adsf", "asd-3", "-"}
		for _, token := range tests {
			fn := func(t *testing.T) {
				newToken := GetToken(token)

				assert.Equal(t, 0, newToken.Limit)
				assert.Equal(t, 0, newToken.Offset)
				assert.Equal(t, newToken.NextToken(10, 100), "")
			}
			t.Run(token, fn)
		}
	})

	t.Run("valid token incremented with records returned equal to limit", func(t *testing.T) {
		tests := []struct {
			limit, offset int
		}{
			{
				limit:  10,
				offset: 0,
			},
			{
				limit:  4,
				offset: 12,
			},
			{
				limit:  3,
				offset: 9,
			},
		}
		for _, tt := range tests {
			fn := func(t *testing.T) {
				tok := Token{
					Limit:  tt.limit,
					Offset: tt.offset,
				}

				var buf bytes.Buffer
				err := gob.NewEncoder(&buf).Encode(tok)
				require.NoError(t, err)

				newToken := GetToken(base64.URLEncoding.EncodeToString(buf.Bytes()))

				assert.Equal(t, tt.limit, newToken.Limit)
				assert.Equal(t, tt.offset, newToken.Offset)

				nextTok := Token{
					Limit:  tt.limit,
					Offset: tt.offset + tt.limit,
				}

				var newBuf bytes.Buffer
				err = gob.NewEncoder(&newBuf).Encode(nextTok)
				require.NoError(t, err)

				expectedToken := base64.URLEncoding.EncodeToString([]byte(newBuf.Bytes()))
				assert.Equal(t, nextTok.String(), expectedToken)
			}
			t.Run(fmt.Sprintf("limit-%d offset-%d", tt.limit, tt.offset), fn)
		}
	})

	t.Run("last page", func(t *testing.T) {
		var tests = []struct {
			limit, offset int
			expected      bool
			total         int
		}{
			{10, 0, true, 100},
			{10, 90, false, 100},
			{11, 90, false, 100},
			{10, 90, true, TotalTODO},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("limit-%d offset-%d", tt.limit, tt.offset), func(t *testing.T) {
				tok := Token{
					Limit:  tt.limit,
					Offset: tt.offset,
				}

				assert.Equal(t, tt.expected, tok.NextToken(10, tt.total) != "")
			})
		}
	})

	t.Run("Params included in token", func(t *testing.T) {
		tests := []struct {
			name   string
			params map[string]string
		}{
			{
				name: "simple1",
				params: map[string]string{
					"param1": "value 1",
					"param2": "value 2",
				},
			},
			{
				name: "simple2",
				params: map[string]string{
					"param1": "value 1",
					"param2": "value 2",
					"param3": "value 3",
					"param4": "value 4",
				},
			},
			{
				name: "empty vals",
				params: map[string]string{
					"param1": "",
					"param2": "",
				},
			},
			{
				name:   "no Params",
				params: map[string]string{},
			},
		}
		for _, tt := range tests {
			fn := func(t *testing.T) {
				tok := Token{
					Offset: 1,
					Limit:  1,
					Params: tt.params,
				}

				newToken := GetToken(tok.String())
				assert.Equal(t, 1, newToken.Limit)
				assert.Equal(t, 1, newToken.Offset)
				require.Equal(t, len(tt.params), len(newToken.Params))

				for key, val := range tt.params {
					assert.Equal(t, val, newToken.GetParam(key))
				}
			}
			t.Run(tt.name, fn)
		}
	})
}

func Test_GetTokenFromQuery(t *testing.T) {
	t.Run("page token not set", func(t *testing.T) {
		t.Run("provided offset sets token offset", func(t *testing.T) {
			q := make(url.Values)
			q.Set("page-token", "")
			q.Set("limit", "5")

			token := GetTokenFromQuery(q)

			assert.Equal(t, 5, token.Limit)
			assert.Zero(t, token.Offset)
		})

		t.Run("provided limit sets token limit", func(t *testing.T) {
			q := make(url.Values)
			q.Set("page-token", "")
			q.Set("limit", "3")
			q.Set("offset", "5")

			token := GetTokenFromQuery(q)

			assert.Equal(t, 3, token.Limit)
			assert.Equal(t, 5, token.Offset)
		})
	})

	t.Run("page token set", func(t *testing.T) {
		t.Run("provided query Params does not override token from page token", func(t *testing.T) {
			prevToken := Token{
				Offset: 30,
				Limit:  30,
			}

			t.Log(prevToken.String())

			q := make(url.Values)
			q.Set("page-token", prevToken.String())
			q.Set("limit", "5")
			q.Set("offset", "5")

			token := GetTokenFromQuery(q)

			assert.Equal(t, prevToken.Limit, token.Limit)
			assert.Equal(t, prevToken.Offset, token.Offset)
		})
	})

	t.Run("DefaultLimit", func(t *testing.T) {
		t.Run("does not set when limit was set in query param", func(t *testing.T) {
			q := make(url.Values)
			q.Set("limit", "5")
			q.Set("offset", "5")

			token := GetTokenFromQuery(q, DefaultLimit(300))

			assert.Equal(t, 5, token.Limit)
			assert.Equal(t, 5, token.Offset)
		})

		t.Run("does not set when pageToken has set limit", func(t *testing.T) {
			prevToken := Token{
				Limit: 30,
			}

			q := make(url.Values)
			q.Set("page-token", prevToken.String())

			token := GetTokenFromQuery(q, DefaultLimit(300))

			assert.Equal(t, prevToken.Limit, token.Limit)
			assert.Equal(t, prevToken.Offset, token.Offset)
		})

		t.Run("does set when no token limit or query limit", func(t *testing.T) {
			q := make(url.Values)
			token := GetTokenFromQuery(q, DefaultLimit(300))

			assert.Equal(t, 300, token.Limit)
			assert.Zero(t, token.Offset)
		})
	})

	t.Run("MinLimit", func(t *testing.T) {
		t.Run("sets min limit when pagetoken limit is less than min", func(t *testing.T) {
			q := make(url.Values)
			q.Set("page-token", Token{Limit: 1}.String())

			pageToken := GetTokenFromQuery(q, MinLimit(30))

			assert.Equal(t, 30, pageToken.Limit)
		})

		t.Run("sets min limit when no pagetoke limit and has query limit less than min", func(t *testing.T) {
			q := make(url.Values)
			q.Set("limit", "0")

			pageToken := GetTokenFromQuery(q, MinLimit(30))

			assert.Equal(t, 30, pageToken.Limit)
		})

		t.Run("does not set when page token limit greater than min", func(t *testing.T) {
			q := make(url.Values)
			q.Set("page-token", Token{Limit: 5}.String())

			pageToken := GetTokenFromQuery(q, MinLimit(1))

			assert.Equal(t, 5, pageToken.Limit)
		})

		t.Run("does not set when query limit greater than min", func(t *testing.T) {
			q := make(url.Values)
			q.Set("limit", "5")

			pageToken := GetTokenFromQuery(q, MinLimit(1))

			assert.Equal(t, 5, pageToken.Limit)
		})
	})

	t.Run("MaxLimit", func(t *testing.T) {
		t.Run("sets max limit when pagetoken limit is greater than max", func(t *testing.T) {
			q := make(url.Values)
			q.Set("page-token", Token{Limit: 100}.String())

			pageToken := GetTokenFromQuery(q, MaxLimit(3))

			assert.Equal(t, 3, pageToken.Limit)
		})

		t.Run("sets max limit when no pagetoke limit and has query limit greater than max", func(t *testing.T) {
			q := make(url.Values)
			q.Set("limit", "10")

			pageToken := GetTokenFromQuery(q, MaxLimit(3))

			assert.Equal(t, 3, pageToken.Limit)
		})

		t.Run("does not set when page token limit less than max", func(t *testing.T) {
			q := make(url.Values)
			q.Set("page-token", Token{Limit: 5}.String())

			pageToken := GetTokenFromQuery(q, MaxLimit(1))

			assert.Equal(t, 1, pageToken.Limit)
		})

		t.Run("does not set when query limit less than max", func(t *testing.T) {
			q := make(url.Values)
			q.Set("limit", "5")

			pageToken := GetTokenFromQuery(q, MaxLimit(1))

			assert.Equal(t, 1, pageToken.Limit)
		})
	})
}

func ExampleToken_first() {
	params := url.Values{
		"limit":  []string{"5"},
		"offset": []string{"0"},
	}

	token := GetTokenFromQuery(params)

	fmt.Println(token.Limit)
	fmt.Println(token.Offset)
	fmt.Println(token.String())
	// Output:
	// 5
	// 0
	// NP-BAwEBBVRva2VuAf-CAAEDAQVMaW1pdAEEAAEGT2Zmc2V0AQQAAQZQYXJhbXMB_4QAAAAh_4MEAQERbWFwW3N0cmluZ11zdHJpbmcB_4QAAQwBDAAABf-CAQoA
}

func ExampleToken_second() {
	params := url.Values{
		TokenQueryParam: []string{"NP-BAwEBBVRva2VuAf-CAAEDAQVMaW1pdAEEAAEGT2Zmc2V0AQQAAQZQYXJhbXMB_4QAAAAh_4MEAQERbWFwW3N0cmluZ11zdHJpbmcB_4QAAQwBDAAABf-CAQoA"},
	}

	token := GetTokenFromQuery(params)

	fmt.Println(token.Limit)
	fmt.Println(token.Offset)
	fmt.Println(token.String())
	// Output:
	// 5
	// 0
	// NP-BAwEBBVRva2VuAf-CAAEDAQVMaW1pdAEEAAEGT2Zmc2V0AQQAAQZQYXJhbXMB_4QAAAAh_4MEAQERbWFwW3N0cmluZ11zdHJpbmcB_4QAAQwBDAAABf-CAQoA
}

func ExampleToken_NextToken_first() {
	params := url.Values{
		"limit":  []string{"5"},
		"offset": []string{"0"},
	}

	token := GetTokenFromQuery(params)

	totalRecords := 6
	nextPageToken := token.NextToken(5, totalRecords)

	fmt.Println(nextPageToken)

	token = GetToken(nextPageToken)
	fmt.Println(token.Limit)
	fmt.Println(token.Offset)
	fmt.Println(token.String())
	// Output:
	// NP-BAwEBBVRva2VuAf-CAAEDAQVMaW1pdAEEAAEGT2Zmc2V0AQQAAQZQYXJhbXMB_4QAAAAh_4MEAQERbWFwW3N0cmluZ11zdHJpbmcB_4QAAQwBDAAAB_-CAQoBCgA=
	// 5
	// 5
	// NP-BAwEBBVRva2VuAf-CAAEDAQVMaW1pdAEEAAEGT2Zmc2V0AQQAAQZQYXJhbXMB_4QAAAAh_4MEAQERbWFwW3N0cmluZ11zdHJpbmcB_4QAAQwBDAAAB_-CAQoBCgA=
}

func ExampleToken_NextToken_second() {
	// strToken below is equivalent to setting limit=5 and offset=0
	strToken := "NP-BAwEBBVRva2VuAf-CAAEDAQVMaW1pdAEEAAEGT2Zmc2V0AQQAAQZQYXJhbXMB_4QAAAAh_4MEAQERbWFwW3N0cmluZ11zdHJpbmcB_4QAAQwBDAAABf-CAQoA"
	params := url.Values{
		TokenQueryParam: []string{strToken},
	}

	token := GetTokenFromQuery(params)

	totalRecords := 6
	nextPageToken := token.NextToken(5, totalRecords)

	fmt.Println(nextPageToken)

	token = GetToken(nextPageToken)
	fmt.Println(token.Limit)
	fmt.Println(token.Offset)
	fmt.Println(token.String())
	// Output:
	// NP-BAwEBBVRva2VuAf-CAAEDAQVMaW1pdAEEAAEGT2Zmc2V0AQQAAQZQYXJhbXMB_4QAAAAh_4MEAQERbWFwW3N0cmluZ11zdHJpbmcB_4QAAQwBDAAAB_-CAQoBCgA=
	// 5
	// 5
	// NP-BAwEBBVRva2VuAf-CAAEDAQVMaW1pdAEEAAEGT2Zmc2V0AQQAAQZQYXJhbXMB_4QAAAAh_4MEAQERbWFwW3N0cmluZ11zdHJpbmcB_4QAAQwBDAAAB_-CAQoBCgA=
}

func ExampleToken_NextToken_third() {
	params := url.Values{
		"limit":  []string{"5"},
		"offset": []string{"0"},
	}

	token := GetTokenFromQuery(params)

	totalRecords := 5
	nextPageToken := token.NextToken(5, totalRecords)
	fmt.Println(nextPageToken)
	// Output:
	//
}

func ExampleToken_NextToken_fourth() {
	// strToken below is equivalent to setting limit=5 and offset=0
	strToken := "NP-BAwEBBVRva2VuAf-CAAEDAQVMaW1pdAEEAAEGT2Zmc2V0AQQAAQZQYXJhbXMB_4QAAAAh_4MEAQERbWFwW3N0cmluZ11zdHJpbmcB_4QAAQwBDAAABf-CAQoA"
	params := url.Values{
		TokenQueryParam: []string{strToken},
	}

	token := GetTokenFromQuery(params)

	totalRecords := 5
	// total records == 5, so no pagination
	nextPageToken := token.NextToken(5, totalRecords)
	fmt.Println(nextPageToken)

	// Output:
	//
}

func ExampleDefaultLimit() {
	params := url.Values{}

	token := GetTokenFromQuery(params, DefaultLimit(10))

	fmt.Println(token.Limit)

	// Output:
	// 10
}

func ExampleMinLimit() {
	params := url.Values{
		"limit": []string{"-2000"},
	}

	token := GetTokenFromQuery(params, MinLimit(1))

	fmt.Println(token.Limit)

	// Output:
	// 1
}

func ExampleMaxLimit() {
	params := url.Values{
		"limit": []string{"2000"},
	}

	token := GetTokenFromQuery(params, MaxLimit(100))

	fmt.Println(token.Limit)

	// Output:
	// 100
}

func ExampleParam() {
	params := url.Values{}

	token := GetTokenFromQuery(params,
		DefaultLimit(10),
		Param("sort", "name"),
		Param("order", "asc"),
	)

	fmt.Println(token.Limit)
	fmt.Println(token.GetParam("sort"))
	fmt.Println(token.GetParam("order"))

	// params are baked into the NextToken
	totalRecords := 11
	nextToken := GetToken(token.NextToken(10, totalRecords))

	fmt.Println(nextToken.Limit)
	fmt.Println(token.GetParam("sort"))
	fmt.Println(token.GetParam("order"))

	// Output:
	// 10
	// name
	// asc
	// 10
	// name
	// asc
}

func ExampleTokenOptFn() {
	params := url.Values{
		"limit": []string{"2000"},
	}

	token := GetTokenFromQuery(params,
		DefaultLimit(10),
		MinLimit(1),
		MaxLimit(100),
		Param("foo", "bar"),
	)

	fmt.Println(token.Limit)

	// Output:
	// 100
}
