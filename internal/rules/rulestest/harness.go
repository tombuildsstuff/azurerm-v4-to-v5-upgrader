package rulestest

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/require"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/engine"
	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

var update = flag.Bool("update", false, "rewrite golden fixture files")

// RunGolden runs r against each case subdirectory of ruleDir.
// Each case dir must contain input.tf; expected.tf and expected-report.json
// are compared (or written when -update is set).
func RunGolden(t *testing.T, r rules.Rule, ruleDir string) {
	t.Helper()
	entries, err := os.ReadDir(ruleDir)
	require.NoError(t, err, "reading %s", ruleDir)

	cases := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		cases++
		caseDir := filepath.Join(ruleDir, e.Name())
		t.Run(e.Name(), func(t *testing.T) {
			runCase(t, []rules.Rule{r}, caseDir)
		})
	}
	require.Positive(t, cases, "rule %s has no fixture cases in %s", r.ID(), ruleDir)
}

// RunGoldenAllRules runs the full default ruleset (rules.All()) against the
// input.tf found directly in caseDir, comparing to expected.tf and
// expected-report.json (or writing them when -update is set). Unlike
// RunGolden, caseDir itself is treated as a single fixture case, not a
// directory of cases.
func RunGoldenAllRules(t *testing.T, caseDir string) {
	t.Helper()
	runCase(t, nil, caseDir)
}

func runCase(t *testing.T, rs []rules.Rule, caseDir string) {
	t.Helper()
	input, err := os.ReadFile(filepath.Join(caseDir, "input.tf"))
	require.NoError(t, err)

	tmp := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tmp, "input.tf"), input, 0o644))

	res, err := engine.Run(engine.Options{Dir: tmp, Rules: rs})
	require.NoError(t, err)

	gotTF, err := os.ReadFile(filepath.Join(tmp, "input.tf"))
	require.NoError(t, err)

	// Guard every fixture's actual transformed output against invalid HCL:
	// golden-file comparison alone can't catch a rule that produces well-
	// formed-looking-but-broken output (e.g. a duplicate attribute) if the
	// expected.tf was generated from that same buggy output via -update.
	_, parseDiags := hclsyntax.ParseConfig(gotTF, "input.tf", hcl.InitialPos)
	require.False(t, parseDiags.HasErrors(), "transformed HCL is invalid: %s", parseDiags.Error())

	var reportBuf bytes.Buffer
	require.NoError(t, res.Report.RenderJSON(&reportBuf))
	// Normalize the random t.TempDir() prefix out so goldens are path-stable.
	normalized := strings.ReplaceAll(reportBuf.String(), tmp+string(os.PathSeparator), "")

	expTFPath := filepath.Join(caseDir, "expected.tf")
	expReportPath := filepath.Join(caseDir, "expected-report.json")

	if *update {
		require.NoError(t, os.WriteFile(expTFPath, gotTF, 0o644))
		require.NoError(t, os.WriteFile(expReportPath, []byte(normalized), 0o644))
		return
	}

	wantTF, err := os.ReadFile(expTFPath)
	require.NoError(t, err, "missing expected.tf (run with -update)")
	require.Equal(t, string(wantTF), string(gotTF), "transformed HCL mismatch")

	wantReport, err := os.ReadFile(expReportPath)
	require.NoError(t, err, "missing expected-report.json (run with -update)")
	require.JSONEq(t, string(wantReport), normalized, "report mismatch")
}
