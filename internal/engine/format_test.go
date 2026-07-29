package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMinimizeFormatChurnRevertsWhitespaceOnlyLines(t *testing.T) {
	// hclwrite canonicalised the colon spacing on an untouched tags line while a
	// real rename happened above it. Only the rename should survive.
	orig := `resource "azurerm_key_vault" "e" {
  enable_rbac_authorization = true
  tags = {
    "source": "terraform"
  }
}
`
	formatted := `resource "azurerm_key_vault" "e" {
  rbac_authorization_enabled = true
  tags = {
    "source" : "terraform"
  }
}
`
	got := string(minimizeFormatChurn([]byte(orig), []byte(formatted), "main.tf"))
	want := `resource "azurerm_key_vault" "e" {
  rbac_authorization_enabled = true
  tags = {
    "source": "terraform"
  }
}
`
	require.Equal(t, want, got)
}

func TestMinimizeFormatChurnKeepsSemanticLineChanges(t *testing.T) {
	// A pure rename with no incidental churn round-trips to the formatted output.
	orig := "resource \"x\" \"y\" {\n  enable_bgp = true\n}\n"
	formatted := "resource \"x\" \"y\" {\n  bgp_enabled = true\n}\n"
	got := string(minimizeFormatChurn([]byte(orig), []byte(formatted), "main.tf"))
	require.Equal(t, formatted, got)
}

func TestMinimizeFormatChurnPreservesAddedAndRemovedLines(t *testing.T) {
	// A block rename that also drops and adds lines: insert/delete hunks are kept
	// verbatim, an untouched whitespace-only line is reverted.
	orig := `resource "x" "y" {
  is_device_condition {
    match_values = ["Mobile"]
  }
  tags = {
    "k": "v"
  }
}
`
	formatted := `resource "x" "y" {
  device_type {
    values = ["Mobile"]
  }
  tags = {
    "k" : "v"
  }
}
`
	got := string(minimizeFormatChurn([]byte(orig), []byte(formatted), "main.tf"))
	require.Contains(t, got, "device_type {")
	require.Contains(t, got, "values = [\"Mobile\"]")
	require.Contains(t, got, `"k": "v"`, "untouched whitespace-only line must be reverted")
	require.NotContains(t, got, `"k" : "v"`)
}

func TestMinimizeFormatChurnFallsBackOnUnparseableResult(t *testing.T) {
	// If reconciliation somehow produced invalid HCL, the safe formatted output
	// is returned. Here orig is not valid HCL, so any reverted line would break
	// parsing; the function must return the (valid) formatted bytes unchanged.
	orig := "this is not : hcl at all\n"
	formatted := "resource \"x\" \"y\" {\n  bgp_enabled = true\n}\n"
	got := string(minimizeFormatChurn([]byte(orig), []byte(formatted), "main.tf"))
	require.Equal(t, formatted, got)
}

func TestMinimizeFormatChurnNoChangeIsIdentity(t *testing.T) {
	src := "resource \"x\" \"y\" {\n  a = 1\n}\n"
	got := string(minimizeFormatChurn([]byte(src), []byte(src), "main.tf"))
	require.Equal(t, src, got)
}

func TestRunPreserveFormattingSuppressesColonChurn(t *testing.T) {
	dir := t.TempDir()
	src := `resource "azurerm_key_vault" "e" {
  enable_rbac_authorization = true
  tags = {
    "source": "terraform"
  }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(src), 0o644))
	_, err := Run(Options{Dir: dir, PreserveFormatting: true})
	require.NoError(t, err)
	out, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)
	got := string(out)
	require.Contains(t, got, "rbac_authorization_enabled = true", "the rename must still be applied")
	require.Contains(t, got, `"source": "terraform"`, "untouched colon line keeps original formatting under --fmt=false")
	require.NotContains(t, got, `"source" : "terraform"`)
}

func TestRunDefaultCanonicallyFormats(t *testing.T) {
	dir := t.TempDir()
	src := "resource \"azurerm_key_vault\" \"e\" {\n  enable_rbac_authorization = true\n  tags = {\n    \"source\": \"terraform\"\n  }\n}\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(src), 0o644))
	_, err := Run(Options{Dir: dir})
	require.NoError(t, err)
	out, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)
	require.Contains(t, string(out), `"source" : "terraform"`, "default output is canonically formatted like terraform fmt")
}
