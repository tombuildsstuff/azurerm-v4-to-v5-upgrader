package engine_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/engine"
	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"

	_ "github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules/azurerm"
)

func copyToTemp(t *testing.T, src string) string {
	t.Helper()
	dir := t.TempDir()
	b, err := os.ReadFile(src)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), b, 0o644))
	return dir
}

func TestNoOpOnPlainConfig(t *testing.T) {
	dir := t.TempDir()
	src := "variable \"name\" {\n  type = string\n}\n\n# a comment\noutput \"o\" {\n  value = var.name\n}\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(src), 0o644))

	// Dry-run must propose no changes for a config with no azurerm resources.
	res, err := engine.Run(engine.Options{Dir: dir, DryRun: true})
	require.NoError(t, err)
	require.Empty(t, res.Diffs)

	// A real run must leave the file byte-identical.
	_, err = engine.Run(engine.Options{Dir: dir})
	require.NoError(t, err)
	got, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)
	require.Equal(t, src, string(got), "untouched config must be byte-identical")
}

func TestIdempotentSecondRun(t *testing.T) {
	dir := copyToTemp(t, filepath.Join("testdata", "corpus", "input", "main.tf"))

	// First real run migrates the config.
	_, err := engine.Run(engine.Options{Dir: dir})
	require.NoError(t, err)
	after1, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)

	// A dry-run over the already-migrated config must propose NO further changes.
	res2, err := engine.Run(engine.Options{Dir: dir, DryRun: true})
	require.NoError(t, err)
	require.Empty(t, res2.Diffs, "second run must propose no changes")

	// A real second run must leave the file byte-identical.
	_, err = engine.Run(engine.Options{Dir: dir})
	require.NoError(t, err)
	after2, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)
	require.Equal(t, string(after1), string(after2))
}

func TestCorpusMatchesExpected(t *testing.T) {
	dir := copyToTemp(t, filepath.Join("testdata", "corpus", "input", "main.tf"))
	res, err := engine.Run(engine.Options{Dir: dir})
	require.NoError(t, err)

	got, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)
	want, err := os.ReadFile(filepath.Join("testdata", "corpus", "expected", "main.tf"))
	require.NoError(t, err)
	require.Equal(t, string(want), string(got))

	// The corpus spans all rule categories (provider pin, top-level and
	// nested value-preserving renames, value-change, name->id ref/literal,
	// changed-default pin, removed block, flag-only inversion, removed
	// resource); assert the structured report reflects that breadth.
	require.GreaterOrEqual(t, res.Report.RulesFired, 11)

	var changed, note, needsReview int
	for _, f := range res.Report.Findings {
		switch f.Severity {
		case rules.SeverityChanged:
			changed++
		case rules.SeverityNote:
			note++
		case rules.SeverityNeedsReview:
			needsReview++
		}
	}
	require.NotZero(t, changed, "expected at least one Changed finding (e.g. a value-preserving rename)")
	require.NotZero(t, note, "expected at least one Note finding (e.g. a pinned changed default)")
	require.NotZero(t, needsReview, "expected at least one NeedsReview finding (e.g. a removed resource)")
	require.Equal(t, 5, needsReview, "azurerm_app_service removal + storage_account value-change/removed-block + storage_queue literal name-to-id + application_insights inversion should all be flagged NeedsReview")
}

func TestCommentsAndInterpolationPreserved(t *testing.T) {
	dir := t.TempDir()
	src := "resource \"azurerm_virtual_network_gateway\" \"gw\" {\n" +
		"  # keep this comment\n" +
		"  name       = \"${var.prefix}-gw\"\n" +
		"  type       = \"Vpn\"\n" +
		"  enable_bgp = true\n" +
		"}\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(src), 0o644))

	_, err := engine.Run(engine.Options{Dir: dir})
	require.NoError(t, err)
	got, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)
	require.Contains(t, string(got), "# keep this comment")
	require.Contains(t, string(got), "\"${var.prefix}-gw\"")
	require.Contains(t, string(got), "bgp_enabled = true")
}
