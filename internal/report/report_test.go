package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

func finding(sev rules.Severity, id string, line int) rules.Finding {
	return rules.Finding{
		RuleID:   id,
		Severity: sev,
		File:     "input.tf",
		Range:    hcl.Range{Start: hcl.Pos{Line: line}},
		Message:  "msg",
	}
}

func TestExitCode(t *testing.T) {
	clean := New()
	clean.Add(finding(rules.SeverityChanged, "a", 1))
	require.Equal(t, 0, clean.ExitCode())

	review := New()
	review.Add(finding(rules.SeverityNeedsReview, "b", 1))
	require.Equal(t, 1, review.ExitCode())
}

func TestRenderJSONDeterministic(t *testing.T) {
	r := New()
	r.ResourcesSeen = 2
	r.RulesFired = 1
	r.Add(finding(rules.SeverityNote, "z", 5))
	r.Add(finding(rules.SeverityChanged, "a", 2))

	var buf bytes.Buffer
	require.NoError(t, r.RenderJSON(&buf))

	var got struct {
		Findings []struct {
			RuleID   string `json:"rule_id"`
			Severity string `json:"severity"`
			Line     int    `json:"line"`
		} `json:"findings"`
		ResourcesSeen int `json:"resources_seen"`
		RulesFired    int `json:"rules_fired"`
	}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
	require.Equal(t, 2, got.ResourcesSeen)
	require.Len(t, got.Findings, 2)
	// Changed (0) sorts before Note (1).
	require.Equal(t, "a", got.Findings[0].RuleID)
	require.Equal(t, "Changed", got.Findings[0].Severity)
	require.Equal(t, "z", got.Findings[1].RuleID)
}

func TestRenderTextIncludesSeverityAndMessage(t *testing.T) {
	r := New()
	r.Add(finding(rules.SeverityNeedsReview, "b", 7))
	var buf bytes.Buffer
	require.NoError(t, r.RenderText(&buf))
	require.Contains(t, buf.String(), "NeedsReview")
	require.Contains(t, buf.String(), "b")
	require.Contains(t, buf.String(), "input.tf:7")
}

func TestRenderMarkdownWithFindings(t *testing.T) {
	r := New()
	r.ResourcesSeen = 2
	r.RulesFired = 1
	r.Add(finding(rules.SeverityNeedsReview, "provider.version-pin", 3))

	var buf bytes.Buffer
	require.NoError(t, r.RenderMarkdown(&buf, true))

	out := buf.String()
	require.Contains(t, out, "## azurerm v5 upgrade check")
	require.Contains(t, out, "Changes required: **yes**")
	require.Contains(t, out, "| NeedsReview | provider.version-pin | input.tf:3 | msg |")
}

func TestRenderMarkdownEscapesTableCells(t *testing.T) {
	r := New()
	f := finding(rules.SeverityNeedsReview, "provider.version-pin", 3)
	f.Message = "a | b\nc"
	r.Add(f)

	var buf bytes.Buffer
	require.NoError(t, r.RenderMarkdown(&buf, true))

	out := buf.String()
	require.Contains(t, out, `a \| b c`)
	require.NotContains(t, out, "| b\n")

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	tableLines := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "| ") {
			tableLines++
		}
	}
	// header row + separator row + exactly one finding row (not split by the
	// embedded newline).
	require.Equal(t, 3, tableLines)
}

func TestRenderMarkdownGroupsByConfigDir(t *testing.T) {
	r := New()
	r.ResourcesSeen = 2
	r.RulesFired = 1
	fA := finding(rules.SeverityNeedsReview, "provider.version-pin", 3)
	fA.File = "configs/network/main.tf"
	fB := finding(rules.SeverityNeedsReview, "provider.version-pin", 4)
	fB.File = "configs/storage/main.tf"
	r.Add(fB) // added out of order to prove directory sorting
	r.Add(fA)

	var buf bytes.Buffer
	require.NoError(t, r.RenderMarkdown(&buf, true))
	out := buf.String()

	require.Contains(t, out, "### configs/network")
	require.Contains(t, out, "### configs/storage")
	require.Less(t, strings.Index(out, "### configs/network"), strings.Index(out, "### configs/storage"),
		"directories must be emitted in sorted order")
	require.Contains(t, out, "configs/network/main.tf:3")
	require.Contains(t, out, "configs/storage/main.tf:4")
}

func TestRenderMarkdownClean(t *testing.T) {
	r := New()
	var buf bytes.Buffer
	require.NoError(t, r.RenderMarkdown(&buf, false))
	require.Contains(t, buf.String(), "Changes required: **no**")
	require.NotContains(t, buf.String(), "| Severity |")
}
