package rules

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func indexOf(s, sub string) int { return strings.Index(s, sub) }

func parseBody(t *testing.T, src string) *hclwrite.Body {
	t.Helper()
	f, diags := hclwrite.ParseConfig([]byte(src), "t.tf", hcl.InitialPos)
	require.False(t, diags.HasErrors(), diags.Error())
	return f.Body()
}

func TestRenameAttributePreservesExpressionCommentAndPosition(t *testing.T) {
	src := "resource \"azurerm_x\" \"a\" {\n" +
		"  old_name = var.thing # keep me\n" +
		"  after    = 1\n" +
		"}\n"
	body := parseBody(t, src)
	blk := body.Blocks()[0]
	require.True(t, RenameAttribute(blk.Body(), "old_name", "new_name"))
	out := string(hclwrite.Format(body.BuildTokens(nil).Bytes()))
	require.Contains(t, out, "new_name")
	require.NotContains(t, out, "old_name")
	require.Contains(t, out, "var.thing")
	require.Contains(t, out, "# keep me") // inline comment preserved
	// position preserved: new_name still appears before `after`
	require.Less(t, indexOf(out, "new_name"), indexOf(out, "after"))
	require.False(t, RenameAttribute(blk.Body(), "missing", "x"))
}

func TestRemoveAttributeReturnsOldText(t *testing.T) {
	body := parseBody(t, "resource \"azurerm_x\" \"a\" {\n  gone = \"value\"\n}\n")
	blk := body.Blocks()[0]
	old, removed := RemoveAttribute(blk.Body(), "gone")
	require.True(t, removed)
	require.Equal(t, "\"value\"", old)
	_, removed = RemoveAttribute(blk.Body(), "gone")
	require.False(t, removed)
}

func TestFindResourceBlocksAndAddress(t *testing.T) {
	body := parseBody(t, "resource \"azurerm_x\" \"a\" {}\nresource \"azurerm_y\" \"b\" {}\n")
	blocks := FindResourceBlocks(body, "azurerm_x")
	require.Len(t, blocks, 1)
	require.Equal(t, "azurerm_x.a", ResourceAddress(blocks[0]))
}

func TestPinChangedDefaultInjectsWhenAbsent(t *testing.T) {
	body := parseBody(t, "resource \"azurerm_x\" \"a\" {\n  name = \"n\"\n}\n")
	blk := body.Blocks()[0]
	require.True(t, PinChangedDefault(blk.Body(), "public_access", ctyStr("Enabled"), "Disabled"))
	// hclwrite only bakes in canonical "=" alignment/spacing when the token
	// stream is formatted; body.BuildTokens alone reflects the raw insertion
	// order without that spacing. Format first, as the other helper tests do,
	// so the assertions are robust to hclwrite's exact spacing choices.
	out := string(hclwrite.Format(body.BuildTokens(nil).Bytes()))
	require.Contains(t, out, "public_access = \"Enabled\"")
	require.Contains(t, out, "# NOTE: AzureRM v5 changes this default to 'Disabled'")
	// The NOTE comment must not introduce a spurious blank line before the
	// closing brace of the block (i.e. the attribute's own newline plus a
	// second newline from the comment token).
	require.NotContains(t, out, "\n\n}")

	// Absent injection is idempotent-safe: present now, so returns false.
	require.False(t, PinChangedDefault(blk.Body(), "public_access", ctyStr("Enabled"), "Disabled"))
}

func TestRenameResourceTypeEmitsMovedBlock(t *testing.T) {
	body := parseBody(t, "resource \"azurerm_old\" \"a\" {\n  name = \"n\"\n}\n")
	blk := body.Blocks()[0]
	RenameResourceType(body, blk, "azurerm_new")
	out := string(hclwrite.Format(body.BuildTokens(nil).Bytes()))
	require.Contains(t, out, "resource \"azurerm_new\" \"a\"")
	require.Contains(t, out, "moved {")
	require.Contains(t, out, "from = azurerm_old.a")
	require.Contains(t, out, "to   = azurerm_new.a")
}

func TestBlockKindStrings(t *testing.T) {
	require.Equal(t, "resource", ResourceBlock.HCLType())
	require.Equal(t, "data", DataBlock.HCLType())
	require.Equal(t, "", ResourceBlock.IDPrefix())
	require.Equal(t, "data.", DataBlock.IDPrefix())
	require.Equal(t, "", ResourceBlock.AddrPrefix())
	require.Equal(t, "data.", DataBlock.AddrPrefix())
}

func TestFindBlocksMatchesKind(t *testing.T) {
	src := "resource \"azurerm_x\" \"a\" {}\n" +
		"data \"azurerm_x\" \"b\" {}\n" +
		"data \"azurerm_y\" \"c\" {}\n"
	body := parseBody(t, src)

	res := FindBlocks(body, ResourceBlock, "azurerm_x")
	require.Len(t, res, 1)
	require.Equal(t, "azurerm_x.a", ResourceAddress(res[0]))

	dat := FindBlocks(body, DataBlock, "azurerm_x")
	require.Len(t, dat, 1)
	require.Equal(t, "azurerm_x.b", ResourceAddress(dat[0]))

	// FindResourceBlocks still resource-only.
	require.Len(t, FindResourceBlocks(body, "azurerm_x"), 1)
	require.Empty(t, FindResourceBlocks(body, "azurerm_y"))
}

func ctyStr(s string) cty.Value { return cty.StringVal(s) }

func TestTraversalRenamesNameToID(t *testing.T) {
	body := parseBody(t, "resource \"azurerm_x\" \"a\" {\n"+
		"  ref     = azurerm_storage_account.example.name\n"+
		"  literal = \"foo\"\n"+
		"  varref  = var.x\n"+
		"  wrong   = azurerm_other.y.name\n"+
		"  attr    = azurerm_storage_account.example.location\n"+
		"}\n")
	blk := body.Blocks()[0]

	attr := blk.Body().GetAttribute("ref")
	tokens, ok := TraversalRenamesNameToID(attr, "azurerm_storage_account")
	require.True(t, ok)
	require.Equal(t, "azurerm_storage_account.example.id", strings.TrimSpace(string(tokens.Bytes())))

	for _, name := range []string{"literal", "varref", "wrong", "attr"} {
		attr := blk.Body().GetAttribute(name)
		_, ok := TraversalRenamesNameToID(attr, "azurerm_storage_account")
		require.False(t, ok, "expected %s to be rejected", name)
	}
}

func TestTraversalRenamesNameToID_DataSource4Part(t *testing.T) {
	body := parseBody(t, "resource \"x\" \"y\" {\n  a = data.azurerm_storage_account.foo.name\n}\n")
	attr := body.Blocks()[0].Body().GetAttribute("a")
	toks, ok := TraversalRenamesNameToID(attr, "azurerm_storage_account")
	require.True(t, ok)
	require.Equal(t, "data.azurerm_storage_account.foo.id", strings.TrimSpace(string(toks.Bytes())))
}

func TestTraversalRenamesNameToID_DataSourceWrongType(t *testing.T) {
	body := parseBody(t, "resource \"x\" \"y\" {\n  a = data.azurerm_other.foo.name\n}\n")
	attr := body.Blocks()[0].Body().GetAttribute("a")
	_, ok := TraversalRenamesNameToID(attr, "azurerm_storage_account")
	require.False(t, ok)
}

func TestTraversalRenamesNameToID_ResourceStillWorks(t *testing.T) {
	body := parseBody(t, "resource \"x\" \"y\" {\n  a = azurerm_storage_account.foo.name\n}\n")
	attr := body.Blocks()[0].Body().GetAttribute("a")
	toks, ok := TraversalRenamesNameToID(attr, "azurerm_storage_account")
	require.True(t, ok)
	require.Equal(t, "azurerm_storage_account.foo.id", strings.TrimSpace(string(toks.Bytes())))
}

func TestWalkNestedBlocks(t *testing.T) {
	src := "resource \"azurerm_x\" \"a\" {\n" +
		"  security {\n    old = true\n  }\n" +
		"  data_disk {\n    v = 1\n  }\n" +
		"  data_disk {\n    v = 2\n  }\n" +
		"}\n"
	body := parseBody(t, src).Blocks()[0].Body()
	require.Len(t, WalkNestedBlocks(body, nil), 1)
	require.Len(t, WalkNestedBlocks(body, []string{"security"}), 1)
	require.Len(t, WalkNestedBlocks(body, []string{"data_disk"}), 2) // repeated
	require.Empty(t, WalkNestedBlocks(body, []string{"missing"}))
}

func TestFindOrCreateBlockReusesThenCreates(t *testing.T) {
	body := parseBody(t, "resource \"azurerm_x\" \"a\" {\n  ip_configuration {\n    name = \"n\"\n  }\n}\n")
	parent := body.Blocks()[0].Body()

	existing := FindOrCreateBlock(parent, "ip_configuration")
	require.NotNil(t, existing)
	require.NotNil(t, existing.Body().GetAttribute("name")) // reused, not new

	created := FindOrCreateBlock(parent, "brand_new")
	require.NotNil(t, created)
	require.Equal(t, "brand_new", created.Type())
	require.Len(t, parent.Blocks(), 2)
}

func TestMoveAttributePreservesCommentsAndReportsAbsence(t *testing.T) {
	src := "resource \"azurerm_x\" \"a\" {\n" +
		"  network_interface {\n" +
		"    # lead comment\n" +
		"    target_static_ip = \"10.0.0.4\" # inline comment\n" +
		"  }\n" +
		"  ip_configuration {\n" +
		"  }\n" +
		"}\n"
	body := parseBody(t, src)
	res := body.Blocks()[0].Body()
	from := res.Blocks()[0].Body() // network_interface
	to := res.Blocks()[1].Body()   // ip_configuration

	require.True(t, MoveAttribute(from, to, "target_static_ip"))
	require.False(t, MoveAttribute(from, to, "target_static_ip")) // gone now

	out := string(hclwrite.Format(body.BuildTokens(nil).Bytes()))
	require.Contains(t, out, "target_static_ip")
	require.Contains(t, out, "# lead comment")   // lead comment moved
	require.Contains(t, out, "# inline comment") // inline comment moved
	// target_static_ip now sits after the ip_configuration opening brace,
	// i.e. after "ip_configuration {".
	require.Greater(t, indexOf(out, "target_static_ip"), indexOf(out, "ip_configuration {"))
}
