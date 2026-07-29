package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/config"
	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

// tagRule sets an attribute on every azurerm_resource_group block.
type tagRule struct{}

func (tagRule) ID() string          { return "test.tag" }
func (tagRule) Description() string { return "tags rg blocks" }
func (tagRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range ctx.Write.Blocks() {
		if b.Type() != "resource" {
			continue
		}
		labels := b.Labels()
		if len(labels) < 1 || labels[0] != "azurerm_resource_group" {
			continue
		}
		b.Body().SetAttributeValue("upgraded", cty.BoolVal(true))
		out = append(out, rules.Finding{RuleID: "test.tag", Severity: rules.SeverityChanged, File: ctx.File})
	}
	return out
}

func TestRunWritesChanges(t *testing.T) {
	dir := t.TempDir()
	src := "resource \"azurerm_resource_group\" \"a\" {\n  name = \"x\"\n}\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(src), 0o644))

	res, err := Run(Options{Dir: dir, Rules: []rules.Rule{tagRule{}}})
	require.NoError(t, err)
	require.Equal(t, 1, res.Report.ResourcesSeen)
	require.Equal(t, 1, res.Report.RulesFired)

	got, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)
	require.Contains(t, string(got), "upgraded = true")
}

func TestRunDryRunDoesNotWrite(t *testing.T) {
	dir := t.TempDir()
	src := "resource \"azurerm_resource_group\" \"a\" {\n  name = \"x\"\n}\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(src), 0o644))

	res, err := Run(Options{Dir: dir, DryRun: true, Rules: []rules.Rule{tagRule{}}})
	require.NoError(t, err)
	require.Contains(t, res.Diffs[filepath.Join(dir, "main.tf")], "upgraded = true")

	got, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)
	require.Equal(t, src, string(got))
}

func TestRunNoChangeLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	src := "variable \"y\" {}\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "vars.tf"), []byte(src), 0o644))

	res, err := Run(Options{Dir: dir, Rules: []rules.Rule{tagRule{}}})
	require.NoError(t, err)
	require.Empty(t, res.Diffs)
	require.Equal(t, 0, res.Report.RulesFired)
}

func TestRunRecursiveLoadsNestedConfigsAndSkipsDotTerraform(t *testing.T) {
	root := t.TempDir()
	write := func(rel, body string) {
		p := filepath.Join(root, rel)
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(body), 0o644))
	}
	write(filepath.Join("net", "main.tf"), `resource "azurerm_virtual_network" "a" {}`)
	write(filepath.Join("stor", "main.tf"), `resource "azurerm_storage_account" "b" {}`)
	write(filepath.Join(".terraform", "modules", "x", "main.tf"), `resource "azurerm_virtual_network" "c" {}`)

	res, err := Run(Options{Dir: root, DryRun: true, Recursive: true})
	require.NoError(t, err)
	require.Equal(t, 2, res.Report.ResourcesSeen, ".terraform config must be skipped")
}

func mustSyntaxBody(t *testing.T, src string) *hclsyntax.Body {
	t.Helper()
	f, diags := hclsyntax.ParseConfig([]byte(src), "t.tf", hcl.InitialPos)
	require.False(t, diags.HasErrors(), diags.Error())
	return f.Body.(*hclsyntax.Body)
}

func TestCountsDataSources(t *testing.T) {
	f := &config.File{
		Syntax: mustSyntaxBody(t, "resource \"azurerm_x\" \"a\" {}\n"+
			"data \"azurerm_y\" \"b\" {}\n"+
			"data \"aws_thing\" \"c\" {}\n"),
	}
	require.Equal(t, 2, countAzureRMResources(f))
}
