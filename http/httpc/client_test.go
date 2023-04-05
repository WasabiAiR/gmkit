package httpc_test

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/graymeta/gmkit/backoff"
	"github.com/graymeta/gmkit/errors"
	"github.com/graymeta/gmkit/http/httpc"
	"github.com/graymeta/gmkit/http/httpc/httpcfakes"
)

func TestClient_Req(t *testing.T) {
	t.Run("no body", func(t *testing.T) {
		t.Run("basics", func(t *testing.T) {
			tests := []int{http.StatusOK, http.StatusAccepted}
			methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

			for _, method := range methods {
				for _, status := range tests {
					fn := func(t *testing.T) {
						doer := new(httpcfakes.FakeDoer)
						doer.DoReturns(stubResp(status), nil)

						client := httpc.New(doer)

						err := client.
							Req(method, "/foo").
							Success(func(statusCode int) bool {
								return status == statusCode
							}).
							Do(context.TODO())
						require.NoError(t, err)
					}

					t.Run(method+"/"+http.StatusText(status), fn)
				}
			}
		})

		t.Run("DELETE", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(http.StatusNoContent), nil)

			client := httpc.New(doer)

			err := client.
				DELETE("/foo").
				Success(httpc.StatusNoContent()).
				Do(context.TODO())
			require.NoError(t, err)
		})
	})

	t.Run("with response headers", func(t *testing.T) {
		svr := httptest.NewServer(testServer(t))
		defer svr.Close()

		expectedHeader := "Yo"
		expectedHeaderValue := "Son"
		doer := new(httpcfakes.FakeDoer)
		doer.DoStub = func(r *http.Request) (*http.Response, error) {
			return stubRespHeaders(http.StatusOK, map[string]string{
				expectedHeader: expectedHeaderValue,
			}), nil
		}

		client := httpc.New(doer, httpc.WithBaseURL(svr.URL))

		t.Run("Do", func(t *testing.T) {
			var totalCalls int
			err := client.
				GET("/foo").
				Success(httpc.StatusOK()).
				ResponseHeaders(func(header http.Header) {
					totalCalls++
					assert.Equal(t, expectedHeaderValue, header.Get(expectedHeader))
				}).
				Do(context.TODO())

			require.NoError(t, err)
			require.Equal(t, 1, totalCalls)
		})

		t.Run("DoAndGetReader", func(t *testing.T) {
			var totalCalls int
			resp, err := client.
				GET("/foo").
				Success(httpc.StatusOK()).
				ResponseHeaders(func(header http.Header) {
					totalCalls++
					assert.Equal(t, expectedHeaderValue, header.Get(expectedHeader))
				}).
				DoAndGetReader(context.TODO())
			defer drain(resp)

			require.NotNil(t, resp)
			require.NoError(t, err)
			require.Equal(t, 1, totalCalls)
		})
	})

	t.Run("with response body", func(t *testing.T) {
		svr := httptest.NewServer(testServer(t))
		defer svr.Close()

		t.Run("GET", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			expected := foo{Name: "Name"}
			doer.DoStub = func(r *http.Request) (*http.Response, error) {
				return stubRespNBody(t, http.StatusOK, expected), nil
			}

			client := httpc.New(doer, httpc.WithBaseURL(svr.URL))

			t.Run("Do Decode", func(t *testing.T) {
				var fooResp foo
				err := client.
					GET("/foo").
					Success(httpc.StatusOK()).
					Decode(httpc.JSONDecode(&fooResp)).
					Do(context.TODO())
				require.NoError(t, err)

				assert.Equal(t, expected, fooResp)
			})

			t.Run("Do DecodeJSON", func(t *testing.T) {
				var fooResp foo
				err := client.
					GET("/foo").
					Success(httpc.StatusOK()).
					DecodeJSON(&fooResp).
					Do(context.TODO())
				require.NoError(t, err)

				assert.Equal(t, expected, fooResp)
			})

			t.Run("DoAndGetReader", func(t *testing.T) {
				var fooResp foo
				resp, err := client.
					GET("/foo").
					Success(httpc.StatusOK()).
					DoAndGetReader(context.TODO())
				drain(resp)

				require.NoError(t, err)
				err = json.NewDecoder(resp.Body).Decode(&fooResp)
				require.NoError(t, err)
				assert.Equal(t, expected, fooResp)
			})
		})

		t.Run("with request body", func(t *testing.T) {
			client := httpc.New(http.DefaultClient, httpc.WithBaseURL(svr.URL))
			expected := foo{Name: "name", S: "string"}

			t.Run("POST Do", func(t *testing.T) {
				var fooResp foo
				err := client.
					POST("/foo").
					Body(expected).
					Success(httpc.StatusOK()).
					Decode(httpc.JSONDecode(&fooResp)).
					Do(context.TODO())

				require.NoError(t, err)
				expected.Method = "POST"
				assert.Equal(t, expected, fooResp)
			})

			t.Run("POST DoAndGetReader", func(t *testing.T) {
				var fooResp foo
				resp, err := client.
					POST("/foo").
					Body(expected).
					Success(httpc.StatusOK()).
					DoAndGetReader(context.TODO())
				defer drain(resp)

				require.NoError(t, err)
				require.NotNil(t, resp)
				expected.Method = "POST"
				err = json.NewDecoder(resp.Body).Decode(&fooResp)
				require.NoError(t, err)
				assert.Equal(t, expected, fooResp)
			})

			t.Run("PATCH Do", func(t *testing.T) {
				var fooResp foo
				err := client.
					PATCH("/foo").
					Body(expected).
					Success(httpc.StatusOK()).
					Decode(httpc.JSONDecode(&fooResp)).
					Do(context.TODO())

				require.NoError(t, err)
				expected.Method = "PATCH"
				assert.Equal(t, expected, fooResp)
			})

			t.Run("PATCH DoAndGetReader", func(t *testing.T) {
				var fooResp foo
				resp, err := client.
					PATCH("/foo").
					Body(expected).
					Success(httpc.StatusOK()).
					DoAndGetReader(context.TODO())
				defer drain(resp)

				require.NoError(t, err)
				require.NotNil(t, resp)
				expected.Method = "PATCH"
				err = json.NewDecoder(resp.Body).Decode(&fooResp)
				require.NoError(t, err)
				assert.Equal(t, expected, fooResp)
			})

			t.Run("PUT Do", func(t *testing.T) {
				var fooResp foo
				err := client.
					PUT("/foo").
					Body(expected).
					Success(httpc.StatusOK()).
					Decode(httpc.JSONDecode(&fooResp)).
					Do(context.TODO())

				require.NoError(t, err)
				expected.Method = "PUT"
				assert.Equal(t, expected, fooResp)
			})

			t.Run("PUT DoAndGetReader", func(t *testing.T) {
				var fooResp foo
				resp, err := client.
					PUT("/foo").
					Body(expected).
					Success(httpc.StatusOK()).
					DoAndGetReader(context.TODO())
				defer drain(resp)

				require.NoError(t, err)
				require.NotNil(t, resp)
				expected.Method = "PUT"
				err = json.NewDecoder(resp.Body).Decode(&fooResp)
				require.NoError(t, err)
				assert.Equal(t, expected, fooResp)
			})
		})
	})

	t.Run("with query params", func(t *testing.T) {
		t.Run("without duplicates", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(200), nil)

			client := httpc.New(doer)

			req := client.GET("/foo").Success(httpc.StatusOK())

			for i := 'A'; i <= 'Z'; i++ {
				req = req.QueryParam(string(i), string(i+26))
			}

			t.Run("Do", func(t *testing.T) {
				err := req.Do(context.TODO())
				require.NoError(t, err)

				httpReq := doer.DoArgsForCall(0)
				params := httpReq.URL.Query()

				for i := 'A'; i <= 'Z'; i++ {
					assert.Equal(t, string(i+26), params.Get(string(i)))
				}
			})

			t.Run("DoAndGetReader", func(t *testing.T) {
				resp, err := req.DoAndGetReader(context.TODO())
				defer drain(resp)
				require.NoError(t, err)
				require.NotNil(t, resp)

				httpReq := doer.DoArgsForCall(0)
				params := httpReq.URL.Query()

				for i := 'A'; i <= 'Z'; i++ {
					assert.Equal(t, string(i+26), params.Get(string(i)))
				}
			})
		})

		t.Run("with duplicates", func(t *testing.T) {
			t.Run("duplicate entries last entry wins", func(t *testing.T) {
				doer := new(httpcfakes.FakeDoer)
				doer.DoReturns(stubResp(200), nil)

				client := httpc.New(doer)

				err := client.
					GET("/foo").
					QueryParam("dupe", "val1").
					QueryParam("dupe", "val2").
					Success(httpc.StatusOK()).
					Do(context.TODO())
				require.NoError(t, err)

				httpReq := doer.DoArgsForCall(0)
				params := httpReq.URL.Query()

				assert.Equal(t, "val2", params.Get("dupe"))
			})
		})

		t.Run("ignores unfulfilled pairs", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(200), nil)

			client := httpc.New(doer)

			err := client.
				GET("/foo").
				QueryParams("q1", "v1", "q2").
				Success(httpc.StatusOK()).
				Do(context.TODO())
			require.NoError(t, err)

			req := doer.DoArgsForCall(0)
			params := req.URL.Query()
			assert.Zero(t, params.Get("q2"))
		})
	})

	t.Run("gob encoding", func(t *testing.T) {
		doer := new(httpcfakes.FakeDoer)
		doer.DoStub = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       r.Body,
			}, nil
		}

		client := httpc.New(doer, httpc.WithEncode(httpc.GobEncode()))
		expected := foo{Name: "name", S: "string"}

		var fooResp foo
		err := client.
			GET("/foo").
			Body(expected).
			Success(httpc.StatusOK()).
			Decode(httpc.GobDecode(&fooResp)).
			Do(context.TODO())

		require.NoError(t, err)
		assert.Equal(t, expected, fooResp)
	})

	t.Run("handling errors response body", func(t *testing.T) {
		type bar struct{ Name string }

		t.Run("Do", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			expected := bar{Name: "error"}
			doer.DoReturns(stubRespNBody(t, http.StatusNotFound, expected), nil)
			client := httpc.New(doer)
			var actual bar
			err := client.
				DELETE("/foo").
				Success(httpc.StatusNoContent()).
				OnError(httpc.JSONDecode(&actual)).
				Do(context.TODO())
			require.Error(t, err)

			// test that verifies the resp body is still readable after err response fn does its read
			found := strings.Contains(err.Error(), `response_body="{\"Name\":\"error\"}\n"`)
			assert.True(t, found)
			assert.Equal(t, expected, actual)
		})

		t.Run("DoAndGetReader", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			expected := bar{Name: "error"}
			doer.DoReturns(stubRespNBody(t, http.StatusNotFound, expected), nil)
			client := httpc.New(doer)
			var actual bar
			resp, err := client.
				DELETE("/foo").
				Success(httpc.StatusNoContent()).
				OnError(httpc.JSONDecode(&actual)).
				DoAndGetReader(context.TODO())
			require.Error(t, err)
			require.Nil(t, resp)

			// test that verifies the resp body is still readable after err response fn does its read
			found := strings.Contains(err.Error(), `response_body="{\"Name\":\"error\"}\n"`)
			assert.True(t, found)
			assert.Equal(t, expected, actual)
		})
	})

	t.Run("retry", func(t *testing.T) {
		t.Run("sets retry", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoStub = func(r *http.Request) (*http.Response, error) {
				return stubResp(http.StatusInternalServerError), nil
			}

			client := httpc.New(doer)

			t.Run("Do", func(t *testing.T) {
				err := client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					Retry(httpc.RetryStatus(httpc.StatusNotIn(http.StatusNoContent, http.StatusNotFound))).
					Do(context.TODO())
				require.Error(t, err)
				isRetryErr(t, err)

				err = client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					RetryStatus(httpc.StatusNotIn(http.StatusNoContent, http.StatusNotFound)).
					Do(context.TODO())
				require.Error(t, err)
				isRetryErr(t, err)

				err = client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					RetryStatusNotIn(http.StatusNoContent, http.StatusNotFound).
					Do(context.TODO())
				require.Error(t, err)
				isRetryErr(t, err)
			})

			t.Run("DoAndGetReader", func(t *testing.T) {
				resp, err := client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					Retry(httpc.RetryStatus(httpc.StatusNotIn(http.StatusNoContent, http.StatusNotFound))).
					DoAndGetReader(context.TODO())
				require.Error(t, err)
				require.Nil(t, resp)
				isRetryErr(t, err)

				resp, err = client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					RetryStatus(httpc.StatusNotIn(http.StatusNoContent, http.StatusNotFound)).
					DoAndGetReader(context.TODO())
				require.Error(t, err)
				require.Nil(t, resp)
				isRetryErr(t, err)

				resp, err = client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					RetryStatusNotIn(http.StatusNoContent, http.StatusNotFound).
					DoAndGetReader(context.TODO())
				require.Error(t, err)
				require.Nil(t, resp)
				isRetryErr(t, err)
			})

		})

		t.Run("does not set retry", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(http.StatusUnprocessableEntity), nil)

			client := httpc.New(doer)

			t.Run("Do", func(t *testing.T) {
				err := client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					Retry(httpc.RetryStatus(httpc.StatusNotIn(http.StatusUnprocessableEntity))).
					Do(context.TODO())
				require.Error(t, err)
				assert.False(t, retryErr(err))
			})

			t.Run("DoAndGetReader", func(t *testing.T) {
				resp, err := client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					Retry(httpc.RetryStatus(httpc.StatusNotIn(http.StatusUnprocessableEntity))).
					DoAndGetReader(context.TODO())
				require.Error(t, err)
				require.Nil(t, resp)
				assert.False(t, retryErr(err))
			})
		})

		t.Run("applies backoff on retry", func(t *testing.T) {
			boffer := backoff.New(
				backoff.InitBackoff(time.Nanosecond),
				backoff.MaxBackoff(time.Nanosecond),
				backoff.MaxCalls(3),
			)

			t.Run("Do", func(t *testing.T) {
				doer := new(httpcfakes.FakeDoer)
				doer.DoReturns(stubResp(http.StatusInternalServerError), nil)
				client := httpc.New(doer, httpc.WithBackoff(boffer))
				err := client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					Retry(httpc.RetryStatus(httpc.StatusNotIn(http.StatusOK))).
					Retry(httpc.RetryStatus(httpc.StatusNotIn(http.StatusNoContent, http.StatusNotFound))).
					Do(context.TODO())
				require.Error(t, err)
				isRetryErr(t, err)
				assert.Equal(t, 3, doer.DoCallCount())
			})

			t.Run("DoAndGetReader", func(t *testing.T) {
				doer := new(httpcfakes.FakeDoer)
				doer.DoReturns(stubResp(http.StatusInternalServerError), nil)
				client := httpc.New(doer, httpc.WithBackoff(boffer))
				resp, err := client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					Retry(httpc.RetryStatus(httpc.StatusNotIn(http.StatusOK))).
					Retry(httpc.RetryStatus(httpc.StatusNotIn(http.StatusNoContent, http.StatusNotFound))).
					DoAndGetReader(context.TODO())
				require.Error(t, err)
				require.Nil(t, resp)
				isRetryErr(t, err)
				assert.Equal(t, 3, doer.DoCallCount())
			})
		})

		t.Run("retries on retryable response error", func(t *testing.T) {
			boffer := backoff.New(
				backoff.InitBackoff(time.Nanosecond),
				backoff.MaxBackoff(time.Nanosecond),
				backoff.MaxCalls(2),
			)

			t.Run("Do", func(t *testing.T) {
				doer := new(httpcfakes.FakeDoer)
				doer.DoReturns(nil, stderrors.New("retry error"))
				client := httpc.New(doer, httpc.WithBackoff(boffer))
				err := client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					Retry(httpc.RetryResponseError(func(e error) error {
						return &fakeRetryErr{e}
					})).
					Do(context.TODO())
				require.Error(t, err)
				isRetryErr(t, err)
				assert.Equal(t, 2, doer.DoCallCount())
			})

			t.Run("DoAndGetReader", func(t *testing.T) {
				doer := new(httpcfakes.FakeDoer)
				doer.DoReturns(nil, stderrors.New("retry error"))
				client := httpc.New(doer, httpc.WithBackoff(boffer))
				resp, err := client.
					DELETE("/foo").
					Success(httpc.StatusNoContent()).
					Retry(httpc.RetryResponseError(func(e error) error {
						return &fakeRetryErr{e}
					})).
					DoAndGetReader(context.TODO())
				require.Error(t, err)
				require.Nil(t, resp)
				isRetryErr(t, err)
				assert.Equal(t, 2, doer.DoCallCount())
			})
		})

		t.Run("retries on http timeouts", func(t *testing.T) {
			t.Run("on request", func(t *testing.T) {
				boffer := backoff.New(
					backoff.InitBackoff(2*time.Nanosecond),
					backoff.MaxBackoff(2*time.Nanosecond),
					backoff.MaxCalls(2),
				)

				doer := &http.Client{
					Timeout: 1 * time.Millisecond,
				}

				client := httpc.New(doer, httpc.WithBackoff(boffer))

				t.Run("Do", func(t *testing.T) {
					err := client.
						GET("https://github.com").
						Success(httpc.StatusOK()).
						Retry(httpc.RetryClientTimeout()).
						Do(context.TODO())
					require.Error(t, err)
					isRetryErr(t, err)
				})

				t.Run("DoAndGetReader", func(t *testing.T) {
					resp, err := client.
						GET("https://github.com").
						Success(httpc.StatusOK()).
						Retry(httpc.RetryClientTimeout()).
						DoAndGetReader(context.TODO())
					require.Error(t, err)
					require.Nil(t, resp)
					isRetryErr(t, err)
				})
			})

			t.Run("on client", func(t *testing.T) {
				boffer := backoff.New(
					backoff.InitBackoff(2*time.Nanosecond),
					backoff.MaxBackoff(2*time.Nanosecond),
					backoff.MaxCalls(2),
				)

				doer := &http.Client{
					Timeout: 1 * time.Millisecond,
				}

				client := httpc.New(doer, httpc.WithBackoff(boffer), httpc.WithRetryClientTimeouts())

				t.Run("Do", func(t *testing.T) {
					err := client.
						GET("https://github.com").
						Success(httpc.StatusOK()).
						Do(context.TODO())
					require.Error(t, err)
					isRetryErr(t, err)
				})

				t.Run("DoAndGetReader", func(t *testing.T) {
					resp, err := client.
						GET("https://github.com").
						Success(httpc.StatusOK()).
						DoAndGetReader(context.TODO())
					require.Error(t, err)
					require.Nil(t, resp)
					isRetryErr(t, err)
				})
			})
		})
	})

	t.Run("response error handled", func(t *testing.T) {
		doer := new(httpcfakes.FakeDoer)
		doer.DoReturns(nil, stderrors.New("unexpected error"))

		client := httpc.New(doer)

		t.Run("Do", func(t *testing.T) {
			var count int
			err := client.
				DELETE("/foo").
				Success(httpc.StatusOK()).
				Retry(httpc.RetryResponseError(func(e error) error {
					count++
					return e
				})).
				Do(context.TODO())
			require.Error(t, err)
			assert.Equal(t, 1, count)
		})

		t.Run("DoAndGetReader", func(t *testing.T) {
			var count int
			resp, err := client.
				DELETE("/foo").
				Success(httpc.StatusOK()).
				Retry(httpc.RetryResponseError(func(e error) error {
					count++
					return e
				})).
				DoAndGetReader(context.TODO())
			require.Error(t, err)
			require.Nil(t, resp)
			assert.Equal(t, 1, count)
		})
	})

	t.Run("headers", func(t *testing.T) {
		t.Run("non duplicates", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(200), nil)

			client := httpc.New(doer)

			req := client.
				GET("/foo")

			for i := 'A'; i <= 'Z'; i++ {
				req = req.Header(string(i), string(i+26))
			}

			t.Run("Do", func(t *testing.T) {
				err := req.Success(httpc.StatusOK()).
					Do(context.TODO())
				require.NoError(t, err)

				httpReq := doer.DoArgsForCall(0)
				headers := httpReq.Header

				for i := 'A'; i <= 'Z'; i++ {
					assert.Equal(t, string(i+26), headers.Get(string(i)))
				}
			})

			t.Run("DoAndGetReader", func(t *testing.T) {
				resp, err := req.Success(httpc.StatusOK()).
					DoAndGetReader(context.TODO())
				defer drain(resp)
				require.NoError(t, err)
				require.NotNil(t, resp)

				httpReq := doer.DoArgsForCall(0)
				headers := httpReq.Header

				for i := 'A'; i <= 'Z'; i++ {
					assert.Equal(t, string(i+26), headers.Get(string(i)))
				}
			})
		})

		t.Run("duplicate entries last entry wins", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(200), nil)

			client := httpc.New(doer)

			err := client.
				GET("/foo").
				Header("dupe", "val1").
				Header("dupe", "val2").
				Success(httpc.StatusOK()).
				Do(context.TODO())
			require.NoError(t, err)

			httpReq := doer.DoArgsForCall(0)
			headers := httpReq.Header

			assert.Equal(t, "val2", headers.Get("dupe"))
		})
	})

	t.Run("meta fields added to errors", func(t *testing.T) {
		doer := new(httpcfakes.FakeDoer)
		doer.DoReturns(stubResp(404), nil)

		client := httpc.New(doer)

		t.Run("Do", func(t *testing.T) {
			err := client.PUT("/foo").
				Success(httpc.StatusOK()).
				Meta("k1", "v1", "k2", "v2", "k3", "v3").
				Do(context.TODO())
			require.Error(t, err)

			assert.Contains(t, err.Error(), `k1="v1"`)
			assert.Contains(t, err.Error(), `k2="v2"`)
			assert.Contains(t, err.Error(), `k3="v3"`)
		})

		t.Run("DoAndGetReader", func(t *testing.T) {
			resp, err := client.PUT("/foo").
				Success(httpc.StatusOK()).
				Meta("k1", "v1", "k2", "v2", "k3", "v3").
				DoAndGetReader(context.TODO())
			require.Error(t, err)
			require.Nil(t, resp)

			assert.Contains(t, err.Error(), `k1="v1"`)
			assert.Contains(t, err.Error(), `k2="v2"`)
			assert.Contains(t, err.Error(), `k3="v3"`)
		})
	})

	t.Run("content type", func(t *testing.T) {
		doer := new(httpcfakes.FakeDoer)
		doer.DoReturns(stubResp(200), nil)

		client := httpc.New(doer)

		t.Run("Do", func(t *testing.T) {
			err := client.
				GET("/foo").
				ContentType("application/json").
				Success(httpc.StatusOK()).
				Do(context.TODO())
			require.NoError(t, err)

			httpReq := doer.DoArgsForCall(0)
			headers := httpReq.Header

			assert.Equal(t, "application/json", headers.Get("Content-Type"))
		})

		t.Run("DoAndGetReader", func(t *testing.T) {
			resp, err := client.
				GET("/foo").
				ContentType("application/json").
				Success(httpc.StatusOK()).
				DoAndGetReader(context.TODO())
			defer drain(resp)
			require.NotNil(t, resp)
			require.NoError(t, err)

			httpReq := doer.DoArgsForCall(0)
			headers := httpReq.Header

			assert.Equal(t, "application/json", headers.Get("Content-Type"))
		})
	})

	t.Run("content length", func(t *testing.T) {
		doer := new(httpcfakes.FakeDoer)
		doer.DoReturns(stubResp(200), nil)

		client := httpc.New(doer)

		buf := strings.NewReader("1234")
		type fakeReader struct {
			*strings.Reader
		}

		t.Run("Do", func(t *testing.T) {
			err := client.
				POST("/foo").
				Body(&fakeReader{Reader: buf}).
				ContentLength(4).
				Success(httpc.StatusOK()).
				Do(context.TODO())
			require.NoError(t, err)
			httpReq := doer.DoArgsForCall(0)
			assert.Equal(t, int64(4), httpReq.ContentLength)
		})

		t.Run("DoAndGetReader", func(t *testing.T) {
			resp, err := client.
				POST("/foo").
				Body(&fakeReader{Reader: buf}).
				ContentLength(4).
				Success(httpc.StatusOK()).
				DoAndGetReader(context.TODO())
			drain(resp)
			require.NoError(t, err)
			require.NotNil(t, resp)
			httpReq := doer.DoArgsForCall(0)
			assert.Equal(t, int64(4), httpReq.ContentLength)
		})
	})

	t.Run("not found", func(t *testing.T) {
		t.Run("Do sets not found", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(http.StatusNotFound), nil)

			client := httpc.New(doer)

			err := client.
				DELETE("/foo").
				Success(httpc.StatusNoContent()).
				NotFound(httpc.StatusNotFound()).
				Do(context.TODO())
			require.Error(t, err)

			assert.True(t, notFoundErr(err))
		})

		t.Run("DoAndGetReader sets not found", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(http.StatusNotFound), nil)

			client := httpc.New(doer)

			resp, err := client.
				DELETE("/foo").
				Success(httpc.StatusNoContent()).
				NotFound(httpc.StatusNotFound()).
				DoAndGetReader(context.TODO())
			require.Error(t, err)
			require.Nil(t, resp)
			assert.True(t, notFoundErr(err))
		})

		t.Run("Do does not set not found", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(http.StatusUnprocessableEntity), nil)

			client := httpc.New(doer)

			err := client.
				DELETE("/foo").
				Success(httpc.StatusNoContent()).
				NotFound(httpc.StatusNotFound()).
				Do(context.TODO())
			require.Error(t, err)

			assert.False(t, notFoundErr(err))
		})

		t.Run("Do does not set not found", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(http.StatusUnprocessableEntity), nil)

			client := httpc.New(doer)

			resp, err := client.
				DELETE("/foo").
				Success(httpc.StatusNoContent()).
				NotFound(httpc.StatusNotFound()).
				DoAndGetReader(context.TODO())
			require.Error(t, err)
			require.Nil(t, resp)

			assert.False(t, notFoundErr(err))
		})
	})

	t.Run("exists", func(t *testing.T) {
		t.Run("Do sets exist", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(http.StatusUnprocessableEntity), nil)

			client := httpc.New(doer)

			err := client.
				DELETE("/foo").
				Success(httpc.StatusNoContent()).
				Exists(httpc.StatusUnprocessableEntity()).
				Do(context.TODO())
			require.Error(t, err)

			assert.True(t, existsErr(err))
		})

		t.Run("DoAndGetReader sets exist", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(http.StatusUnprocessableEntity), nil)

			client := httpc.New(doer)

			resp, err := client.
				DELETE("/foo").
				Success(httpc.StatusNoContent()).
				Exists(httpc.StatusUnprocessableEntity()).
				DoAndGetReader(context.TODO())
			require.Error(t, err)
			require.Nil(t, resp)

			assert.True(t, existsErr(err))
		})

		t.Run("Do does not set exist", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(http.StatusNotFound), nil)

			client := httpc.New(doer)

			err := client.
				DELETE("/foo").
				Success(httpc.StatusNoContent()).
				Exists(httpc.StatusUnprocessableEntity()).
				Do(context.TODO())
			require.Error(t, err)

			assert.False(t, existsErr(err))
		})

		t.Run("DoAndGetReader does not set exist", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoReturns(stubResp(http.StatusNotFound), nil)

			client := httpc.New(doer)

			resp, err := client.
				DELETE("/foo").
				Success(httpc.StatusNoContent()).
				Exists(httpc.StatusUnprocessableEntity()).
				DoAndGetReader(context.TODO())
			require.Error(t, err)
			require.Nil(t, resp)

			assert.False(t, existsErr(err))
		})
	})

	t.Run("auth", func(t *testing.T) {
		t.Run("Do basic auth", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoStub = func(r *http.Request) (*http.Response, error) {
				u, p, ok := r.BasicAuth()
				if !ok {
					return stubResp(http.StatusInternalServerError), nil
				}
				f := foo{Name: u, S: p, Method: r.Method}
				return stubRespNBody(t, http.StatusOK, f), nil
			}

			client := httpc.New(doer, httpc.WithAuth(httpc.BasicAuth("user", "pass")))

			var actual foo
			err := client.
				GET("/foo").
				Success(httpc.StatusOK()).
				Decode(httpc.JSONDecode(&actual)).
				Do(context.TODO())
			require.NoError(t, err)

			assert.Equal(t, "user", actual.Name)
			assert.Equal(t, "pass", actual.S)
		})

		t.Run("DoAndGetReader basic auth", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoStub = func(r *http.Request) (*http.Response, error) {
				u, p, ok := r.BasicAuth()
				if !ok {
					return stubResp(http.StatusInternalServerError), nil
				}
				f := foo{Name: u, S: p, Method: r.Method}
				return stubRespNBody(t, http.StatusOK, f), nil
			}

			client := httpc.New(doer, httpc.WithAuth(httpc.BasicAuth("user", "pass")))

			var actual foo
			resp, err := client.
				GET("/foo").
				Success(httpc.StatusOK()).
				Decode(httpc.JSONDecode(&actual)).
				DoAndGetReader(context.TODO())
			defer drain(resp)
			require.NoError(t, err)
			require.NotNil(t, resp)

			err = json.NewDecoder(resp.Body).Decode(&actual)
			require.NoError(t, err)
			assert.Equal(t, "user", actual.Name)
			assert.Equal(t, "pass", actual.S)
		})

		t.Run("Do token auth", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoStub = func(r *http.Request) (*http.Response, error) {
				f := foo{Name: r.Header.Get("Authorization"), Method: r.Method}
				return stubRespNBody(t, http.StatusOK, f), nil
			}

			client := httpc.New(doer, httpc.WithAuth(httpc.BearerTokenAuth("token")))

			var actual foo
			err := client.
				GET("/foo").
				Success(httpc.StatusOK()).
				Decode(httpc.JSONDecode(&actual)).
				Do(context.TODO())
			require.NoError(t, err)

			assert.Equal(t, "Bearer token", actual.Name)
		})

		t.Run("Do token auth", func(t *testing.T) {
			doer := new(httpcfakes.FakeDoer)
			doer.DoStub = func(r *http.Request) (*http.Response, error) {
				f := foo{Name: r.Header.Get("Authorization"), Method: r.Method}
				return stubRespNBody(t, http.StatusOK, f), nil
			}

			client := httpc.New(doer, httpc.WithAuth(httpc.BearerTokenAuth("token")))

			var actual foo
			resp, err := client.
				GET("/foo").
				Success(httpc.StatusOK()).
				Decode(httpc.JSONDecode(&actual)).
				DoAndGetReader(context.TODO())
			defer drain(resp)
			require.NoError(t, err)
			require.NotNil(t, resp)

			err = json.NewDecoder(resp.Body).Decode(&actual)
			require.NoError(t, err)
			assert.Equal(t, "Bearer token", actual.Name)
		})
	})
}

type foo struct {
	Name   string
	S      string
	Method string
}

func stubRespNBody(t *testing.T, status int, v any) *http.Response {
	t.Helper()

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		t.Fatal(err)
	}
	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(&buf),
	}
}

func stubRespHeaders(status int, headers map[string]string) *http.Response {
	respHeader := http.Header(map[string][]string{})
	for k, v := range headers {
		respHeader.Set(k, v)
	}

	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
		Header:     respHeader,
	}
}

func stubResp(status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(new(bytes.Buffer)),
	}
}

func testServer(t *testing.T) http.Handler {
	t.Helper()

	r := http.NewServeMux()
	r.HandleFunc("/foo", func(w http.ResponseWriter, req *http.Request) {
		var f foo
		if err := json.NewDecoder(req.Body).Decode(&f); err != nil {
			t.Log("error in decoder: ", err)
			return
		}

		f.Method = req.Method
		if err := json.NewEncoder(w).Encode(f); err != nil {
			t.Log("error in encoder: ", err)
			return
		}
	})

	return r
}

func isRetryErr(t *testing.T, err error) {
	t.Helper()
	if !retryErr(err) {
		t.Fatal("got a non retriable error: ", err)
	}
}

func retryErr(err error) bool {
	r, ok := err.(errors.Retrier)
	return ok && r.Retry()
}

func notFoundErr(err error) bool {
	nf, ok := err.(errors.NotFounder)
	return ok && nf.NotFound()
}

func existsErr(err error) bool {
	ex, ok := err.(errors.Exister)
	return ok && ex.Exists()
}

type fakeRetryErr struct {
	error
}

func (f *fakeRetryErr) Retry() bool {
	return true
}

func drain(resp *http.Response) {
	if resp != nil {
		resp.Body.Close()
	}
}
