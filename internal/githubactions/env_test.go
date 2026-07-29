package githubactions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookupEnvDefaultsAPIURL(t *testing.T) {
	env := LookupEnv(func(k string) string {
		return map[string]string{
			"GITHUB_ACTIONS":    "true",
			"GITHUB_REPOSITORY": "owner/repo",
			"GITHUB_EVENT_NAME": "pull_request",
			"GITHUB_TOKEN":      "t0ken",
		}[k]
	})
	require.True(t, env.IsActions)
	require.Equal(t, "owner/repo", env.Repository)
	require.True(t, env.IsPullRequest())
	require.Equal(t, "https://api.github.com", env.APIURL)
}

func TestLookupEnvHonorsCustomAPIURL(t *testing.T) {
	env := LookupEnv(func(k string) string {
		if k == "GITHUB_API_URL" {
			return "https://ghe.example.com/api/v3"
		}
		return ""
	})
	require.Equal(t, "https://ghe.example.com/api/v3", env.APIURL)
	require.False(t, env.IsActions)
	require.False(t, env.IsPullRequest())
}

func TestParsePRNumberFromPullRequest(t *testing.T) {
	n, err := ParsePRNumber([]byte(`{"pull_request":{"number":42}}`))
	require.NoError(t, err)
	require.Equal(t, 42, n)
}

func TestParsePRNumberFallsBackToTopLevel(t *testing.T) {
	n, err := ParsePRNumber([]byte(`{"number":7}`))
	require.NoError(t, err)
	require.Equal(t, 7, n)
}

func TestParsePRNumberMissing(t *testing.T) {
	_, err := ParsePRNumber([]byte(`{}`))
	require.Error(t, err)
}
