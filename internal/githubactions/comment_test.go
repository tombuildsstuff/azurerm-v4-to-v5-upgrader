package githubactions

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestPostCommentSuccess(t *testing.T) {
	var gotReq *http.Request
	var gotBody []byte
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		gotReq = r
		gotBody, _ = io.ReadAll(r.Body)
		return &http.Response{StatusCode: 201, Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})}

	env := Environment{Repository: "owner/repo", Token: "t0ken", APIURL: "https://api.github.com"}
	require.NoError(t, PostComment(client, env, 42, "hello **world**"))

	require.Equal(t, http.MethodPost, gotReq.Method)
	require.Equal(t, "https://api.github.com/repos/owner/repo/issues/42/comments", gotReq.URL.String())
	require.Equal(t, "Bearer t0ken", gotReq.Header.Get("Authorization"))
	require.Equal(t, "application/vnd.github+json", gotReq.Header.Get("Accept"))

	var payload struct {
		Body string `json:"body"`
	}
	require.NoError(t, json.Unmarshal(gotBody, &payload))
	require.Equal(t, "hello **world**", payload.Body)
}

func TestPostCommentHonorsCustomAPIURLWithTrailingSlash(t *testing.T) {
	var url string
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		url = r.URL.String()
		return &http.Response{StatusCode: 201, Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})}
	env := Environment{Repository: "o/r", Token: "t", APIURL: "https://ghe.example.com/api/v3/"}
	require.NoError(t, PostComment(client, env, 3, "x"))
	require.Equal(t, "https://ghe.example.com/api/v3/repos/o/r/issues/3/comments", url)
}

func TestPostCommentNon2xxIsError(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 403, Body: io.NopCloser(strings.NewReader(`{"message":"forbidden"}`))}, nil
	})}
	env := Environment{Repository: "o/r", Token: "t", APIURL: "https://api.github.com"}
	err := PostComment(client, env, 1, "x")
	require.Error(t, err)
	require.Contains(t, err.Error(), "403")
}
