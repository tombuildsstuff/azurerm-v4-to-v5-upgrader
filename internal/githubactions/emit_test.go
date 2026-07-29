package githubactions

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

func okClient(captured *string) *http.Client {
	return &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if captured != nil {
			*captured = r.URL.String()
		}
		return &http.Response{StatusCode: 201, Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})}
}

func sampleFindings() []rules.Finding {
	return []rules.Finding{{
		RuleID: "provider.version-pin", Severity: rules.SeverityNeedsReview,
		File: "main.tf", Range: hcl.Range{Start: hcl.Pos{Line: 3}}, Message: "pin it",
	}}
}

func TestEmitFullPullRequestContext(t *testing.T) {
	var stdout, summary, output bytes.Buffer
	var postedURL string
	d := Deps{
		Env: Environment{
			IsActions: true, EventName: "pull_request", Repository: "o/r",
			Token: "t", APIURL: "https://api.github.com", EventPath: "/evt",
		},
		Stdout: &stdout, Summary: &summary, Output: &output,
		ReadFile:   func(string) ([]byte, error) { return []byte(`{"pull_request":{"number":9}}`), nil },
		HTTPClient: okClient(&postedURL),
	}
	require.NoError(t, Emit(d, sampleFindings(), "MD BODY", true))

	require.Contains(t, stdout.String(), "::error file=main.tf,line=3::pin it")
	require.Equal(t, "MD BODY", summary.String())
	require.Contains(t, output.String(), "needs-changes=true")
	require.Equal(t, "https://api.github.com/repos/o/r/issues/9/comments", postedURL)
}

func TestEmitPullRequestMissingTokenErrors(t *testing.T) {
	d := Deps{
		Env:    Environment{IsActions: true, EventName: "pull_request", Repository: "o/r"},
		Stdout: &bytes.Buffer{},
	}
	err := Emit(d, sampleFindings(), "MD", true)
	require.Error(t, err)
	require.Contains(t, err.Error(), "GITHUB_TOKEN")
}

func TestEmitNonPullRequestSkipsComment(t *testing.T) {
	var stdout, output bytes.Buffer
	posted := false
	d := Deps{
		Env:        Environment{IsActions: true, EventName: "push"},
		Stdout:     &stdout,
		Output:     &output,
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { posted = true; return nil, nil })},
	}
	require.NoError(t, Emit(d, sampleFindings(), "MD", true))
	require.False(t, posted, "no comment outside a PR context")
	require.Contains(t, output.String(), "needs-changes=true")
}

func TestEmitNotInActionsSkipsAnnotations(t *testing.T) {
	var stdout bytes.Buffer
	d := Deps{Env: Environment{IsActions: false}, Stdout: &stdout}
	require.NoError(t, Emit(d, sampleFindings(), "MD", true))
	require.Empty(t, stdout.String(), "annotations only emitted inside Actions")
}
