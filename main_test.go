package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/githubactions"
)

const v4Config = `terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
  }
}

resource "azurerm_virtual_network_gateway" "gw" {
  name       = "gw"
  type       = "Vpn"
  enable_bgp = true
}
`

const v5Config = `terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=5.0.0"
    }
  }
}
`

func writeConfig(t *testing.T, contents string) string {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(contents), 0o644))
	return dir
}

func TestRunValidateV4ConfigNeedsUpgrade(t *testing.T) {
	dir := writeConfig(t, v4Config)
	before, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)

	var buf bytes.Buffer
	needsUpgrade, err := runValidate(&buf, dir, "json", false, githubactions.Environment{}, nil)
	require.NoError(t, err)
	require.True(t, needsUpgrade, "v4 config with enable_bgp should need an upgrade")

	var got struct {
		NeedsUpgrade bool `json:"needs_upgrade"`
	}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
	require.True(t, got.NeedsUpgrade)

	after, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)
	require.Equal(t, string(before), string(after), "validate must never write files")
}

func TestRunValidateV5ConfigNoUpgradeNeeded(t *testing.T) {
	dir := writeConfig(t, v5Config)
	before, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)

	var buf bytes.Buffer
	needsUpgrade, err := runValidate(&buf, dir, "json", false, githubactions.Environment{}, nil)
	require.NoError(t, err)
	require.False(t, needsUpgrade, "already-v5 config should not need an upgrade")

	var got struct {
		NeedsUpgrade bool `json:"needs_upgrade"`
	}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
	require.False(t, got.NeedsUpgrade)

	after, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)
	require.Equal(t, string(before), string(after), "validate must never write files")
}

func TestRunValidateTextReportIncludesNeedsUpgradeLine(t *testing.T) {
	dir := writeConfig(t, v4Config)

	var buf bytes.Buffer
	needsUpgrade, err := runValidate(&buf, dir, "text", false, githubactions.Environment{}, nil)
	require.NoError(t, err)
	require.True(t, needsUpgrade)
	require.Contains(t, buf.String(), "needs upgrade: true")
}

func TestRunValidateInvalidReportKind(t *testing.T) {
	dir := writeConfig(t, v5Config)

	var buf bytes.Buffer
	_, err := runValidate(&buf, dir, "yaml", false, githubactions.Environment{}, nil)
	require.Error(t, err)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestRunValidateGitHubActionsV4NeedsChangesAndPosts(t *testing.T) {
	dir := writeConfig(t, v4Config)
	eventPath := filepath.Join(t.TempDir(), "event.json")
	require.NoError(t, os.WriteFile(eventPath, []byte(`{"pull_request":{"number":5}}`), 0o644))
	summaryPath := filepath.Join(t.TempDir(), "summary.md")
	outputPath := filepath.Join(t.TempDir(), "output.txt")

	var postedURL, postedBody string
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		postedURL = r.URL.String()
		b, _ := io.ReadAll(r.Body)
		postedBody = string(b)
		return &http.Response{StatusCode: 201, Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})}

	env := githubactions.Environment{
		IsActions: true, EventName: "pull_request", Repository: "o/r", Token: "t",
		APIURL: "https://api.github.com", EventPath: eventPath,
		SummaryPath: summaryPath, OutputPath: outputPath,
	}

	var stdout bytes.Buffer
	needs, err := runValidate(&stdout, dir, "text", true, env, client)
	require.NoError(t, err)
	require.True(t, needs)

	require.Equal(t, "https://api.github.com/repos/o/r/issues/5/comments", postedURL)
	require.Contains(t, postedBody, "Changes required")
	summary, _ := os.ReadFile(summaryPath)
	require.Contains(t, string(summary), "## azurerm v5 upgrade check")
	output, _ := os.ReadFile(outputPath)
	require.Contains(t, string(output), "needs-changes=true")
}

func TestRunValidateJSONReportWithGitHubActions(t *testing.T) {
	// --report=json and --github-actions are orthogonal: JSON goes to stdout,
	// grouped markdown goes to the step summary.
	dir := writeConfig(t, v4Config)
	summaryPath := filepath.Join(t.TempDir(), "summary.md")
	env := githubactions.Environment{IsActions: true, EventName: "push", SummaryPath: summaryPath}

	var stdout bytes.Buffer
	needs, err := runValidate(&stdout, dir, "json", true, env, nil)
	require.NoError(t, err)
	require.True(t, needs)

	// stdout leads with the JSON report; GitHub annotations follow it on the
	// same stream (a workflow-command line per finding), so decode the leading
	// JSON value rather than requiring the whole buffer to be pure JSON.
	var got struct {
		NeedsUpgrade bool `json:"needs_upgrade"`
	}
	require.NoError(t, json.NewDecoder(&stdout).Decode(&got), "stdout must lead with the JSON report")
	require.True(t, got.NeedsUpgrade)

	summary, err := os.ReadFile(summaryPath)
	require.NoError(t, err)
	require.Contains(t, string(summary), "## azurerm v5 upgrade check", "step summary must get the markdown")
}

func TestRunValidateGitHubActionsCleanConfigNoComment(t *testing.T) {
	dir := writeConfig(t, v5Config)
	posted := false
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		posted = true
		return &http.Response{StatusCode: 201, Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})}
	env := githubactions.Environment{IsActions: true, EventName: "push"}

	var stdout bytes.Buffer
	needs, err := runValidate(&stdout, dir, "text", true, env, client)
	require.NoError(t, err)
	require.False(t, needs)
	require.False(t, posted)
}

func TestRunValidateGitHubActionsPullRequestWithoutTokenErrors(t *testing.T) {
	dir := writeConfig(t, v4Config)
	env := githubactions.Environment{IsActions: true, EventName: "pull_request", Repository: "o/r"}
	var stdout bytes.Buffer
	_, err := runValidate(&stdout, dir, "text", true, env, http.DefaultClient)
	require.Error(t, err)
	require.Contains(t, err.Error(), "GITHUB_TOKEN")
}

func TestRunUpdateRecursivelyUpgradesNestedConfigs(t *testing.T) {
	root := t.TempDir()
	write := func(rel, body string) {
		p := filepath.Join(root, rel)
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(body), 0o644))
	}
	nested := filepath.Join("modules", "network", "main.tf")
	write(nested, v4Config)                                 // needs upgrade, in a subdir
	write(filepath.Join(".terraform", "main.tf"), v4Config) // must be ignored

	// The root directory itself has no *.tf files - mirrors the reported bug
	// where `validate` saw nested resources but `update` saw zero.
	var buf bytes.Buffer
	_, err := runUpdate(&buf, root, false, true, false, "text")
	require.NoError(t, err)

	// The nested file was rewritten in place: enable_bgp -> bgp_enabled.
	got, err := os.ReadFile(filepath.Join(root, nested))
	require.NoError(t, err)
	require.Contains(t, string(got), "bgp_enabled", "nested config must be upgraded")
	require.NotContains(t, string(got), "enable_bgp", "v4 field must be gone from nested config")

	// .terraform is excluded from the recursive walk.
	skipped, err := os.ReadFile(filepath.Join(root, ".terraform", "main.tf"))
	require.NoError(t, err)
	require.Contains(t, string(skipped), "enable_bgp", ".terraform must not be touched")
}

func TestRunUpdateNonRecursiveOnlyTopLevel(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "sub"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "sub", "main.tf"), []byte(v4Config), 0o644))

	var buf bytes.Buffer
	_, err := runUpdate(&buf, root, false, false, false, "text")
	require.NoError(t, err)

	got, err := os.ReadFile(filepath.Join(root, "sub", "main.tf"))
	require.NoError(t, err)
	require.Contains(t, string(got), "enable_bgp", "non-recursive update must not descend into subdirs")
}

func TestNormalizeArgs(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{"single-dash long flags become double-dash", []string{"update", "-fmt=false", "-dry-run"}, []string{"update", "--fmt=false", "--dry-run"}},
		{"double-dash flags untouched", []string{"update", "--fmt=false"}, []string{"update", "--fmt=false"}},
		{"single-char shorthand untouched", []string{"-h"}, []string{"-h"}},
		{"single-char with value untouched", []string{"-f=x"}, []string{"-f=x"}},
		{"flag value with separate arg", []string{"update", "-dir", "modules"}, []string{"update", "--dir", "modules"}},
		{"bare dash untouched", []string{"-"}, []string{"-"}},
		{"everything after terminator untouched", []string{"-fmt", "--", "-fmt"}, []string{"--fmt", "--", "-fmt"}},
		{"positional args untouched", []string{"validate", "modules"}, []string{"validate", "modules"}},
		{"empty", nil, []string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, normalizeArgs(tc.in))
		})
	}
}

func TestUpdateCmdAcceptsSingleDashFlags(t *testing.T) {
	cmd := updateCmd()
	require.NoError(t, cmd.ParseFlags(normalizeArgs([]string{"-fmt=false", "-dry-run", "-dir", "modules", "-recursive=false"})))

	format, err := cmd.Flags().GetBool("fmt")
	require.NoError(t, err)
	require.False(t, format)

	dryRun, err := cmd.Flags().GetBool("dry-run")
	require.NoError(t, err)
	require.True(t, dryRun)

	dir, err := cmd.Flags().GetString("dir")
	require.NoError(t, err)
	require.Equal(t, "modules", dir)

	recursive, err := cmd.Flags().GetBool("recursive")
	require.NoError(t, err)
	require.False(t, recursive)
}

func TestValidateCmdAcceptsSingleDashFlags(t *testing.T) {
	cmd := validateCmd()
	require.NoError(t, cmd.ParseFlags(normalizeArgs([]string{"-report=json", "-github-actions"})))

	report, err := cmd.Flags().GetString("report")
	require.NoError(t, err)
	require.Equal(t, "json", report)

	gha, err := cmd.Flags().GetBool("github-actions")
	require.NoError(t, err)
	require.True(t, gha)
}

func TestRunValidateRecursivelyDiscoversConfigs(t *testing.T) {
	root := t.TempDir()
	write := func(rel, body string) {
		p := filepath.Join(root, rel)
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(body), 0o644))
	}
	write(filepath.Join("dirty", "main.tf"), v4Config)      // needs upgrade
	write(filepath.Join("clean", "main.tf"), v5Config)      // already v5
	write(filepath.Join(".terraform", "main.tf"), v4Config) // must be ignored

	var buf bytes.Buffer
	needsUpgrade, err := runValidate(&buf, root, "text", false, githubactions.Environment{}, nil)
	require.NoError(t, err)
	require.True(t, needsUpgrade, "a dirty nested config must make the whole run need an upgrade")

	out := buf.String()
	require.Contains(t, out, filepath.Join("dirty", "main.tf"), "dirty config finding must be reported")
	require.NotContains(t, out, filepath.Join(".terraform", "main.tf"), ".terraform must be excluded")
}
