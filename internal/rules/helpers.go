package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// BlockKind distinguishes the two top-level block types a rule can target.
// ResourceBlock is the zero value, so a rule that omits the field targets
// `resource` blocks exactly as before.
type BlockKind int

const (
	ResourceBlock BlockKind = iota // `resource "<type>" "<name>"`
	DataBlock                      // `data "<type>" "<name>"`
)

// HCLType is the block keyword this kind matches.
func (k BlockKind) HCLType() string {
	if k == DataBlock {
		return "data"
	}
	return "resource"
}

// IDPrefix disambiguates a data-source rule's ID from the same-named resource
// rule's ID. Empty for resources so existing IDs are unchanged.
func (k BlockKind) IDPrefix() string {
	if k == DataBlock {
		return "data."
	}
	return ""
}

// AddrPrefix prefixes a finding's Resource address so a data-source finding
// reads `data.azurerm_x.name`. Empty for resources.
func (k BlockKind) AddrPrefix() string {
	if k == DataBlock {
		return "data."
	}
	return ""
}

// FindBlocks returns every `<kind> "<typeLabel>" "*"` block in body.
func FindBlocks(body *hclwrite.Body, kind BlockKind, typeLabel string) []*hclwrite.Block {
	var out []*hclwrite.Block
	for _, b := range body.Blocks() {
		if b.Type() != kind.HCLType() {
			continue
		}
		labels := b.Labels()
		if len(labels) >= 1 && labels[0] == typeLabel {
			out = append(out, b)
		}
	}
	return out
}

// RenameAttribute renames old to new in place: position, expression, lead
// comments, and inline (line) comments are all preserved - only the name
// identifier token changes. Returns false if old is absent (or if new
// already exists on body, matching hclwrite's own conflict handling).
//
// This delegates to (*hclwrite.Body).RenameAttribute, which swaps just the
// attribute's name node via Attribute.setName/node.ReplaceWith rather than
// rebuilding the attribute from expression tokens, so the surrounding
// comment and newline nodes - and the attribute's position among its
// siblings - are left untouched.
func RenameAttribute(body *hclwrite.Body, old, new string) bool {
	return body.RenameAttribute(old, new)
}

// RenameBlockType replaces blk's block-type keyword with newType, preserving
// blk's position among parent's other block-typed children, its labels, its
// body content verbatim (order, comments, formatting), and its lead comment
// (any comment written directly above the block) via a token copy. The old
// blk is left detached from parent; callers must not continue to use it and
// should operate on the returned block instead.
//
// This works around a bug in hashicorp/hcl/v2 v2.24.0 (the latest release as
// of writing): (*hclwrite.Block).SetType writes the new type name into the
// block's serialized token stream (so the file's bytes ARE correct) but,
// unlike every other rename method in that package - Attribute.setName,
// Body.SetAttributeRaw/SetAttributeValue/SetAttributeTraversal, which all
// correctly reassign their node field from node.ReplaceWith's return value -
// SetType discards that return value and leaves its own typeName field
// pointing at the old, now-detached node. So Block.Type() on that same block
// object keeps reporting the OLD type name for the rest of the process, even
// though the emitted file shows the new one. A later WalkNestedBlocks lookup
// keyed on the new type name then silently fails to match, so any rule
// chained after a renamedBlockRule that needs to re-find the renamed block
// (rather than just renaming attributes inside it in the same pass) cannot
// use SetType. Reproduced independently against vanilla hclwrite v2.24.0
// with no project code involved. If a future hcl/v2 upgrade fixes SetType,
// this workaround becomes unnecessary (though still correct).
//
// The replacement block is produced by re-parsing blk's own complete token
// stream (hclwrite.Block.BuildTokens), which includes its lead comment, type
// identifier, labels, and body - with only the leading type identifier token
// rewritten to newType - rather than by rebuilding from the body alone. A
// body-only rebuild would drop any comment written directly above the block,
// since a block's lead comments belong to the Block node, not its Body.
// Re-parsing the rewritten stream (instead of e.g. hclwrite.NewBlock +
// AppendUnstructuredTokens, which would carry the old body's attributes and
// nested blocks over as an opaque, unstructured token blob - the tradeoff
// MoveAttribute/fieldMoveToNestedBlockRule deliberately accept for attributes
// that are never queried again) yields fully structured *Attribute and
// *Block nodes, in original order, comments intact, still queryable via
// GetAttribute/Blocks/WalkNestedBlocks for whatever rule runs next.
//
// Limitation: reconstructing block order this way only preserves order
// *among parent's block-typed children* - any top-level attributes in parent
// are left exactly where they are, so the renamed block (and its re-appended
// sibling blocks) end up ordered after parent's attributes even if the
// original block was interleaved among them. Every current call site (an
// actions{}/conditions{} body containing only nested blocks, no sibling
// attributes) is unaffected by this, since there are no attributes to be
// ordered relative to.
//
// A second consequence of the remove-then-reappend approach: it moves only
// the *Block nodes themselves, not the free-floating newline/comment nodes
// that sit between them in parent's token stream. A standalone blank line or
// comment positioned BETWEEN two sibling blocks is left behind at its
// original position rather than travelling with either neighbor, so after
// the rename the (re-appended) blocks end up contiguous and that blank
// line/comment resurfaces above the renamed block instead of between the
// blocks it used to separate. The output is still valid, semantically
// equivalent HCL - this is a formatting-only change. A block's own lead
// comment (written directly above it, with no intervening blank line) is
// unaffected: that's part of the block's own token stream (see above) and
// moves with it correctly.
// TODO(v5): a position-preserving re-append (splicing the new block's tokens
// in place of the old one, rather than remove-all/re-append-all) would avoid
// this, but is a riskier rewrite than justified for the current call sites.
func RenameBlockType(parent *hclwrite.Body, blk *hclwrite.Block, newType string) (*hclwrite.Block, bool) {
	siblings := parent.Blocks()

	// blk's full token stream: lead comments, type ident, labels, body. Only
	// the first TokenIdent - the block's type name - is rewritten; lead
	// comments are whole-line TokenComment tokens and labels are quoted
	// TokenOQuote/TokenQuotedLit/TokenCQuote strings, so neither can be
	// mistaken for the type identifier.
	toks := blk.BuildTokens(nil)
	var renamed bool
	out := make(hclwrite.Tokens, len(toks))
	for i, t := range toks {
		if !renamed && t.Type == hclsyntax.TokenIdent {
			nt := *t
			nt.Bytes = []byte(newType)
			out[i] = &nt
			renamed = true
			continue
		}
		out[i] = t
	}

	tmpFile, diags := hclwrite.ParseConfig(out.Bytes(), "<rename-block-type>", hcl.InitialPos)
	if diags.HasErrors() {
		// The source was built from an already-successfully-parsed block's
		// own tokens, so a parse failure here shouldn't be reachable in
		// practice. But this tool runs over arbitrary user files, and a
		// crash of the whole upgrade is worse than one skipped rename:
		// degrade by leaving blk (and the rename) untouched, and tell the
		// caller so it doesn't report a rename that didn't happen.
		return blk, false
	}
	newBlk := tmpFile.Body().Blocks()[0]
	tmpFile.Body().RemoveBlock(newBlk)

	for _, s := range siblings {
		parent.RemoveBlock(s)
	}
	for _, s := range siblings {
		if s == blk {
			parent.AppendBlock(newBlk)
		} else {
			parent.AppendBlock(s)
		}
	}
	return newBlk, true
}

// RemoveAttribute deletes name and returns the trimmed text of its old
// expression. removed is false if name was absent.
func RemoveAttribute(body *hclwrite.Body, name string) (oldText string, removed bool) {
	attr := body.GetAttribute(name)
	if attr == nil {
		return "", false
	}
	oldText = strings.TrimSpace(string(attr.Expr().BuildTokens(nil).Bytes()))
	body.RemoveAttribute(name)
	return oldText, true
}

// FindOrCreateBlock returns the first unlabeled child block of blockType,
// appending a new empty one if none exists.
func FindOrCreateBlock(body *hclwrite.Body, blockType string) *hclwrite.Block {
	if b := body.FirstMatchingBlock(blockType, nil); b != nil {
		return b
	}
	return body.AppendNewBlock(blockType, nil)
}

// MoveAttribute moves the attribute named name from body from to body to,
// preserving its lead and inline comments. It works by detaching the whole
// attribute node (RemoveAttribute returns it) and re-appending its complete
// token stream - lead comments, name, '=', expression, inline comment, and
// trailing newline - into to. hclwrite's formatter re-indents the appended
// tokens to to's nesting level on render.
//
// Consequence: the moved attribute lands in to as unstructured tokens, not a
// structured *Attribute, so to.GetAttribute(name) will NOT find it afterward.
// Callers must not rely on querying a moved attribute in its new home.
//
// Returns false (no-op) when name is absent from from.
func MoveAttribute(from, to *hclwrite.Body, name string) bool {
	removed := from.RemoveAttribute(name)
	if removed == nil {
		return false
	}
	to.AppendUnstructuredTokens(removed.BuildTokens(nil))
	return true
}

// FindResourceBlocks returns every `resource "<typeLabel>" "*"` block in body.
func FindResourceBlocks(body *hclwrite.Body, typeLabel string) []*hclwrite.Block {
	return FindBlocks(body, ResourceBlock, typeLabel)
}

// WalkNestedBlocks returns the bodies of all nested blocks reachable from
// body by following path, one block-type name per level. An empty path
// returns []*hclwrite.Body{body}. Repeated blocks at any level (e.g. multiple
// `data_disk` blocks) are all included, and their descendants are all
// followed for the remainder of path.
func WalkNestedBlocks(body *hclwrite.Body, path []string) []*hclwrite.Body {
	current := []*hclwrite.Body{body}
	for _, name := range path {
		var next []*hclwrite.Body
		for _, b := range current {
			for _, blk := range b.Blocks() {
				if blk.Type() == name {
					next = append(next, blk.Body())
				}
			}
		}
		current = next
	}
	return current
}

// ResourceAddress returns "<type>.<name>" for a resource block.
func ResourceAddress(b *hclwrite.Block) string {
	labels := b.Labels()
	if len(labels) < 2 {
		return ""
	}
	return labels[0] + "." + labels[1]
}

// PinChangedDefault, when name is absent, writes name = <oldDefault> with a
// trailing NOTE comment recording the v5 default change. Returns false if the
// attribute is already present (leave the user's explicit value alone).
func PinChangedDefault(body *hclwrite.Body, name string, oldDefault cty.Value, newDefaultDesc string) bool {
	if body.GetAttribute(name) != nil {
		return false
	}
	valueTokens := hclwrite.TokensForValue(oldDefault)
	// Deliberately no trailing "\n" in Bytes: hclwrite.Body.SetAttributeRaw
	// already appends its own newline token after the attribute's token
	// stream. A comment token whose Bytes end in "\n" would terminate the
	// line itself, so hclwrite's newline would land on its own empty line,
	// rendering as a spurious blank line before the block's closing brace.
	// hclsyntax.TokenComment for a "#" comment is still required to end at a
	// newline per the HCL spec, but that requirement is satisfied by the
	// newline hclwrite adds after the attribute - we must not add a second
	// one.
	comment := &hclwrite.Token{
		Type:         hclsyntax.TokenComment,
		Bytes:        []byte(fmt.Sprintf("# NOTE: AzureRM v5 changes this default to '%s'", newDefaultDesc)),
		SpacesBefore: 1,
	}
	body.SetAttributeRaw(name, append(valueTokens, comment))
	return true
}

// TraversalRenamesNameToID inspects attr's expression: if it is a bare
// 3-step scope traversal shaped <targetType>.<label>.name (e.g.
// azurerm_storage_account.example.name), or a 4-step data-source traversal
// shaped data.<targetType>.<label>.name (e.g.
// data.azurerm_storage_account.example.name), it returns tokens for the same
// traversal with its final step rewritten to id - i.e.
// azurerm_storage_account.example.id or
// data.azurerm_storage_account.example.id - and ok=true.
//
// Any other expression shape is rejected with ok=false: a literal string, a
// variable/local reference, an interpolation, a traversal of a different
// resource type, or a traversal that doesn't end in .name. Those cases
// require a human to pick the replacement resource ID, so the caller must
// leave the config untouched.
func TraversalRenamesNameToID(attr *hclwrite.Attribute, targetType string) (hclwrite.Tokens, bool) {
	exprTokens := attr.Expr().BuildTokens(nil)
	parsed, diags := hclsyntax.ParseExpression(exprTokens.Bytes(), "", hcl.InitialPos)
	if diags.HasErrors() {
		return nil, false
	}
	trav, ok := parsed.(*hclsyntax.ScopeTraversalExpr)
	if !ok {
		return nil, false
	}
	// Accept a 3-part resource ref (<targetType>.<label>.name) or a 4-part
	// data-source ref (data.<targetType>.<label>.name); both migrate the
	// trailing .name to .id.
	var lastStep hcl.Traverser
	switch len(trav.Traversal) {
	case 3:
		root, ok := trav.Traversal[0].(hcl.TraverseRoot)
		if !ok || root.Name != targetType {
			return nil, false
		}
		if _, ok := trav.Traversal[1].(hcl.TraverseAttr); !ok {
			return nil, false
		}
		lastStep = trav.Traversal[2]
	case 4:
		root, ok := trav.Traversal[0].(hcl.TraverseRoot)
		if !ok || root.Name != "data" {
			return nil, false
		}
		typeAttr, ok := trav.Traversal[1].(hcl.TraverseAttr)
		if !ok || typeAttr.Name != targetType {
			return nil, false
		}
		if _, ok := trav.Traversal[2].(hcl.TraverseAttr); !ok {
			return nil, false
		}
		lastStep = trav.Traversal[3]
	default:
		return nil, false
	}
	last, ok := lastStep.(hcl.TraverseAttr)
	if !ok || last.Name != "name" {
		return nil, false
	}

	// Rebuild from a copy of the expression's own tokens (never mutate the
	// tokens returned by BuildTokens in place - they're shared with the
	// live document tree) and rewrite the final identifier's bytes from
	// "name" to "id".
	out := make(hclwrite.Tokens, len(exprTokens))
	copy(out, exprTokens)
	lastIdentIdx := -1
	for i, tok := range out {
		if tok.Type == hclsyntax.TokenIdent {
			lastIdentIdx = i
		}
	}
	if lastIdentIdx == -1 {
		return nil, false
	}
	rewritten := *out[lastIdentIdx]
	rewritten.Bytes = []byte("id")
	out[lastIdentIdx] = &rewritten
	return out, true
}

// RenameResourceType rewrites block's first label to newType and appends a
// top-level moved block mapping the old address to the new one.
func RenameResourceType(body *hclwrite.Body, block *hclwrite.Block, newType string) {
	labels := block.Labels()
	if len(labels) < 2 {
		return
	}
	oldType, name := labels[0], labels[1]
	block.SetLabels([]string{newType, name})

	body.AppendNewline()
	moved := body.AppendNewBlock("moved", nil)
	moved.Body().SetAttributeTraversal("from", hcl.Traversal{
		hcl.TraverseRoot{Name: oldType},
		hcl.TraverseAttr{Name: name},
	})
	moved.Body().SetAttributeTraversal("to", hcl.Traversal{
		hcl.TraverseRoot{Name: newType},
		hcl.TraverseAttr{Name: name},
	})
}
