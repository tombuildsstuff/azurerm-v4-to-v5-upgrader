package attr

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

// renamedAttrRule renames oldAttr → newAttr on every block of resourceType.
// When blockPath is non-empty, the rename is applied within each nested
// block reachable by that path (e.g. []string{"security"} for a rename
// scoped to `security { ... }` sub-blocks) rather than at the resource's
// top level.
type renamedAttrRule struct {
	blockKind    rules.BlockKind
	resourceType string
	blockPath    []string
	oldAttr      string
	newAttr      string
}

func (r renamedAttrRule) ID() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("%s%s.%s.renamed", r.blockKind.IDPrefix(), r.resourceType, r.oldAttr)
	}
	return fmt.Sprintf("%s%s.%s.%s.renamed", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.blockPath, "."), r.oldAttr)
}
func (r renamedAttrRule) Description() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("Rename %s.%s to %s in v5", r.resourceType, r.oldAttr, r.newAttr)
	}
	return fmt.Sprintf("Rename %s.%s.%s to %s in v5", r.resourceType, strings.Join(r.blockPath, "."), r.oldAttr, r.newAttr)
}
func (r renamedAttrRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		for _, nb := range rules.WalkNestedBlocks(b.Body(), r.blockPath) {
			if rules.RenameAttribute(nb, r.oldAttr, r.newAttr) {
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityChanged,
					Resource: r.blockKind.AddrPrefix() + rules.ResourceAddress(b),
					File:     ctx.File,
					Message:  fmt.Sprintf("renamed %s to %s", r.oldAttr, r.newAttr),
				})
			}
		}
	}
	return out
}

// nameToIDRule renames oldAttr → newAttr on every block of resourceType,
// where v5 changes the attribute's meaning from a resource NAME to a
// resource ID (e.g. storage_account_name → storage_account_id). This is
// only safe to auto-migrate when the value is itself a reference to the
// matching resource's .name (a traversal shaped
// <targetType>.<label>.name) - the rename then also rewrites that
// traversal's final step to .id. Any other value (a literal, a variable
// reference, an interpolation, or a traversal of the wrong shape/type) is
// left byte-identical and flagged for human review, since there's no way
// to derive the resource ID from it automatically.
type nameToIDRule struct {
	blockKind    rules.BlockKind
	resourceType string
	blockPath    []string
	oldAttr      string
	newAttr      string
	targetType   string
}

func (r nameToIDRule) ID() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("%s%s.%s.name-to-id", r.blockKind.IDPrefix(), r.resourceType, r.oldAttr)
	}
	return fmt.Sprintf("%s%s.%s.%s.name-to-id", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.blockPath, "."), r.oldAttr)
}
func (r nameToIDRule) Description() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("Migrate %s.%s to %s in v5 (name -> id)", r.resourceType, r.oldAttr, r.newAttr)
	}
	return fmt.Sprintf("Migrate %s.%s.%s to %s in v5 (name -> id)", r.resourceType, strings.Join(r.blockPath, "."), r.oldAttr, r.newAttr)
}
func (r nameToIDRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		for _, nb := range rules.WalkNestedBlocks(b.Body(), r.blockPath) {
			attr := nb.GetAttribute(r.oldAttr)
			if attr == nil {
				continue
			}
			idTokens, ok := rules.TraversalRenamesNameToID(attr, r.targetType)
			if !ok {
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityNeedsReview,
					Resource: r.blockKind.AddrPrefix() + rules.ResourceAddress(b),
					File:     ctx.File,
					Message:  fmt.Sprintf("%s was renamed to %s in v5, and its value must change from a name to a resource ID; update manually", r.oldAttr, r.newAttr),
				})
				continue
			}
			rules.RenameAttribute(nb, r.oldAttr, r.newAttr)
			nb.SetAttributeRaw(r.newAttr, idTokens)
			out = append(out, rules.Finding{
				RuleID:   r.ID(),
				Severity: rules.SeverityChanged,
				Resource: r.blockKind.AddrPrefix() + rules.ResourceAddress(b),
				File:     ctx.File,
				Message:  fmt.Sprintf("renamed %s to %s and rewrote its value from .name to .id", r.oldAttr, r.newAttr),
			})
		}
	}
	return out
}

// removedAttrRule removes attr from every block of resourceType (or, when
// blockPath is non-empty, from each nested block reachable by that path)
// and flags it.
type removedAttrRule struct {
	blockKind    rules.BlockKind
	resourceType string
	blockPath    []string
	attr         string
}

func (r removedAttrRule) ID() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("%s%s.%s.removed", r.blockKind.IDPrefix(), r.resourceType, r.attr)
	}
	return fmt.Sprintf("%s%s.%s.%s.removed", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.blockPath, "."), r.attr)
}
func (r removedAttrRule) Description() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("Remove %s.%s in v5 (needs review)", r.resourceType, r.attr)
	}
	return fmt.Sprintf("Remove %s.%s.%s in v5 (needs review)", r.resourceType, strings.Join(r.blockPath, "."), r.attr)
}
func (r removedAttrRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		for _, nb := range rules.WalkNestedBlocks(b.Body(), r.blockPath) {
			if old, removed := rules.RemoveAttribute(nb, r.attr); removed {
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityNeedsReview,
					Resource: r.blockKind.AddrPrefix() + rules.ResourceAddress(b),
					File:     ctx.File,
					Message:  fmt.Sprintf("removed %s (was %s); confirm no behavior change", r.attr, old),
				})
			}
		}
	}
	return out
}

// valueChangeRule flags an attribute whose value is a literal string no
// longer accepted in v5 (removedValues). It never edits the config - since
// there's no single correct replacement to auto-pick, it only surfaces a
// NeedsReview finding for a human to resolve. Non-literal expressions (var
// references, interpolations, function calls) and literal values not in
// removedValues are left untouched with no finding.
type valueChangeRule struct {
	blockKind     rules.BlockKind
	resourceType  string
	blockPath     []string
	attr          string
	removedValues []string
}

func (r valueChangeRule) ID() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("%s%s.%s.value-removed", r.blockKind.IDPrefix(), r.resourceType, r.attr)
	}
	return fmt.Sprintf("%s%s.%s.%s.value-removed", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.blockPath, "."), r.attr)
}
func (r valueChangeRule) Description() string {
	values := strings.Join(r.removedValues, ", ")
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("Flag %s.%s values (%s) no longer accepted in v5 (needs review)", r.resourceType, r.attr, values)
	}
	return fmt.Sprintf("Flag %s.%s.%s values (%s) no longer accepted in v5 (needs review)", r.resourceType, strings.Join(r.blockPath, "."), r.attr, values)
}
func (r valueChangeRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		for _, nb := range rules.WalkNestedBlocks(b.Body(), r.blockPath) {
			attr := nb.GetAttribute(r.attr)
			if attr == nil {
				continue
			}
			value, ok := literalStringValue(attr)
			if !ok || !containsString(r.removedValues, value) {
				continue
			}
			out = append(out, rules.Finding{
				RuleID:   r.ID(),
				Severity: rules.SeverityNeedsReview,
				Resource: r.blockKind.AddrPrefix() + rules.ResourceAddress(b),
				File:     ctx.File,
				Message:  fmt.Sprintf("%s is set to %q, which is no longer accepted in v5; choose a replacement value", r.attr, value),
			})
		}
	}
	return out
}

// flagAttrRule flags a present attribute as needing manual review, without
// making any edit. Used for v5 changes too complex to auto-migrate - e.g.
// multi-attribute consolidations or id<->url swaps - where there's no single
// mechanical rewrite the tool can safely apply.
type flagAttrRule struct {
	blockKind    rules.BlockKind
	resourceType string
	blockPath    []string
	attr         string
	reason       string // human explanation, appended to the message
}

func (r flagAttrRule) ID() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("%s%s.%s.review", r.blockKind.IDPrefix(), r.resourceType, r.attr)
	}
	return fmt.Sprintf("%s%s.%s.%s.review", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.blockPath, "."), r.attr)
}
func (r flagAttrRule) Description() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("Flag %s.%s for manual review in v5 (%s)", r.resourceType, r.attr, r.reason)
	}
	return fmt.Sprintf("Flag %s.%s.%s for manual review in v5 (%s)", r.resourceType, strings.Join(r.blockPath, "."), r.attr, r.reason)
}
func (r flagAttrRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		for _, nb := range rules.WalkNestedBlocks(b.Body(), r.blockPath) {
			if nb.GetAttribute(r.attr) == nil {
				continue
			}
			out = append(out, rules.Finding{
				RuleID:   r.ID(),
				Severity: rules.SeverityNeedsReview,
				Resource: r.blockKind.AddrPrefix() + rules.ResourceAddress(b),
				File:     ctx.File,
				Message:  fmt.Sprintf("%s requires manual review in v5: %s", r.attr, r.reason),
			})
		}
	}
	return out
}

// removedBlockRule flags a nested block that v5 removes entirely as needing
// manual review, without making any edit. Used when the replacement isn't a
// simple rename but a structural change the tool can't safely apply on its
// own (e.g. a nested block that must become a separate resource).
type removedBlockRule struct {
	blockKind    rules.BlockKind
	resourceType string
	blockPath    []string // path to (and including) the removed block
	reason       string
}

func (r removedBlockRule) ID() string {
	return fmt.Sprintf("%s%s.%s.block-removed", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.blockPath, "."))
}
func (r removedBlockRule) Description() string {
	return fmt.Sprintf("Flag %s.%s for manual review in v5 (block removed: %s)", r.resourceType, strings.Join(r.blockPath, "."), r.reason)
}
func (r removedBlockRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	blockName := strings.Join(r.blockPath, ".")
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		if len(rules.WalkNestedBlocks(b.Body(), r.blockPath)) == 0 {
			continue
		}
		out = append(out, rules.Finding{
			RuleID:   r.ID(),
			Severity: rules.SeverityNeedsReview,
			Resource: r.blockKind.AddrPrefix() + rules.ResourceAddress(b),
			File:     ctx.File,
			Message:  fmt.Sprintf("%s block requires manual review in v5: %s", blockName, r.reason),
		})
	}
	return out
}

// requiredAttrRule flags an attribute that was optional pre-v5 but becomes
// required in v5, without making any edit. The tool has no safe value to
// supply, so it only surfaces a NeedsReview finding when the attribute is
// ABSENT; when present (any value), nothing is flagged.
type requiredAttrRule struct {
	blockKind    rules.BlockKind
	resourceType string
	blockPath    []string
	attr         string
	reason       string
}

func (r requiredAttrRule) ID() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("%s%s.%s.now-required", r.blockKind.IDPrefix(), r.resourceType, r.attr)
	}
	return fmt.Sprintf("%s%s.%s.%s.now-required", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.blockPath, "."), r.attr)
}
func (r requiredAttrRule) Description() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("Flag %s.%s as now-required in v5 when absent (%s)", r.resourceType, r.attr, r.reason)
	}
	return fmt.Sprintf("Flag %s.%s.%s as now-required in v5 when absent (%s)", r.resourceType, strings.Join(r.blockPath, "."), r.attr, r.reason)
}
func (r requiredAttrRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		for _, nb := range rules.WalkNestedBlocks(b.Body(), r.blockPath) {
			if nb.GetAttribute(r.attr) != nil {
				continue
			}
			out = append(out, rules.Finding{
				RuleID:   r.ID(),
				Severity: rules.SeverityNeedsReview,
				Resource: r.blockKind.AddrPrefix() + rules.ResourceAddress(b),
				File:     ctx.File,
				Message:  fmt.Sprintf("%s is now required in v5: %s", r.attr, r.reason),
			})
		}
	}
	return out
}

// requiredBlockRule flags a nested block that was optional pre-v5 but
// becomes required in v5, without making any edit. The tool has no safe
// value to supply, so it only surfaces a single NeedsReview finding per
// resource when the block is ABSENT; when present, nothing is flagged.
type requiredBlockRule struct {
	blockKind    rules.BlockKind
	resourceType string
	blockPath    []string // path to the now-required block
	reason       string
}

func (r requiredBlockRule) ID() string {
	return fmt.Sprintf("%s%s.%s.block-now-required", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.blockPath, "."))
}
func (r requiredBlockRule) Description() string {
	return fmt.Sprintf("Flag %s.%s as now-required in v5 when absent (%s)", r.resourceType, strings.Join(r.blockPath, "."), r.reason)
}
func (r requiredBlockRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	blockName := strings.Join(r.blockPath, ".")
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		if len(rules.WalkNestedBlocks(b.Body(), r.blockPath)) > 0 {
			continue
		}
		out = append(out, rules.Finding{
			RuleID:   r.ID(),
			Severity: rules.SeverityNeedsReview,
			Resource: r.blockKind.AddrPrefix() + rules.ResourceAddress(b),
			File:     ctx.File,
			Message:  fmt.Sprintf("the %s block is now required in v5: %s", blockName, r.reason),
		})
	}
	return out
}

// renamedBlockRule renames a nested block's type from oldBlock to newBlock on
// every matched block, preserving the block's labels, attributes, nested
// blocks, and comments (via (*Block).SetType). When parentPath is non-empty,
// the target block is sought inside each nested block reachable by that path.
// After the rename, any innerRenames are applied to the renamed block's body
// (attr rename) and any innerRemoves are deleted from it. Value-preserving, so
// each renamed block yields one SeverityChanged finding. Absent oldBlock is a
// no-op, which also makes the rule idempotent (a re-run sees only newBlock).
type renamedBlockRule struct {
	blockKind    rules.BlockKind
	resourceType string
	parentPath   []string
	oldBlock     string
	newBlock     string
	// innerRenames pairs must be independent: the set of source names (map
	// keys) and target names (map values) must be disjoint. The map is
	// ranged in nondeterministic order, so an overlapping/chained rename
	// (e.g. a target name that is also another entry's source name) would
	// produce order-dependent output.
	innerRenames map[string]string
	innerRemoves []string
}

func (r renamedBlockRule) ID() string {
	if len(r.parentPath) == 0 {
		return fmt.Sprintf("%s%s.%s.block-renamed", r.blockKind.IDPrefix(), r.resourceType, r.oldBlock)
	}
	return fmt.Sprintf("%s%s.%s.%s.block-renamed", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.parentPath, "."), r.oldBlock)
}

func (r renamedBlockRule) Description() string {
	if len(r.parentPath) == 0 {
		return fmt.Sprintf("Rename block %s.%s to %s in v5", r.resourceType, r.oldBlock, r.newBlock)
	}
	return fmt.Sprintf("Rename block %s.%s.%s to %s in v5", r.resourceType, strings.Join(r.parentPath, "."), r.oldBlock, r.newBlock)
}

func (r renamedBlockRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		for _, parent := range rules.WalkNestedBlocks(b.Body(), r.parentPath) {
			for _, blk := range parent.Blocks() {
				if blk.Type() != r.oldBlock {
					continue
				}
				newBlk, ok := rules.RenameBlockType(parent, blk, r.newBlock)
				if !ok {
					// Re-parse failure: RenameBlockType left blk (and the
					// rename) untouched, so there's nothing to apply inner
					// renames/removes to and no rename to report.
					continue
				}
				for oldA, newA := range r.innerRenames {
					rules.RenameAttribute(newBlk.Body(), oldA, newA)
				}
				for _, a := range r.innerRemoves {
					rules.RemoveAttribute(newBlk.Body(), a)
				}
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityChanged,
					Resource: r.blockKind.AddrPrefix() + rules.ResourceAddress(b),
					File:     ctx.File,
					Message:  fmt.Sprintf("renamed block %s to %s", r.oldBlock, r.newBlock),
				})
			}
		}
	}
	return out
}

// requiredChildAttr names a child-block attribute that becomes required in v5,
// with a human-readable reason for the review finding.
type requiredChildAttr struct {
	name   string
	reason string
}

// fieldMoveToNestedBlockRule moves a set of attributes out of a parent block
// and into a child block (childBlock), creating the child when absent. Moves
// preserve comments (via rules.MoveAttribute), so a moved attribute lands in
// the child as unstructured tokens - do not query moved attrs via GetAttribute.
// If the child block already defines an attribute being moved (e.g. a
// partially-migrated or hand-written config that sets it in both places),
// moving it unconditionally would append a second argument of the same name -
// invalid HCL ("Attribute redefined"). Such attrs are left in place on the
// parent and flagged SeverityNeedsReview instead of being moved; the
// SeverityChanged "moved ..." finding lists only the attrs actually moved,
// and is omitted entirely when none were.
//
// For each requiredChildAttrs entry absent from an existing/created child, a
// SeverityNeedsReview finding is emitted (the tool cannot invent the value);
// required attrs are only checked when the child block exists, avoiding noise
// on unrelated parent blocks. If a parent already contains multiple unlabeled
// childBlock blocks, FindOrCreateBlock/FirstMatchingBlock operate on the
// first one only - moved attrs land there, and only the first is checked for
// required attrs - which matches the intended single-child shapes. Each name
// in requiredChildAttrs must NOT also appear in attrs: a moved attr lands as
// unstructured tokens that GetAttribute cannot see, so a name that is both
// moved and required would be falsely flagged as absent.
type fieldMoveToNestedBlockRule struct {
	blockKind          rules.BlockKind
	resourceType       string
	parentPath         []string
	childBlock         string
	attrs              []string
	requiredChildAttrs []requiredChildAttr
	// requiredChildAttrsWhenMultiple names attrs that only become required
	// once a parent has more than one childBlock - e.g. an ip_configuration's
	// name is only required for disambiguation when a network_interface has
	// multiple ip_configuration blocks. Checked independently of
	// requiredChildAttrs: every matching childBlock (not just the first) is
	// checked, and only when the parent's childBlock count exceeds 1.
	requiredChildAttrsWhenMultiple []requiredChildAttr
}

func (r fieldMoveToNestedBlockRule) ID() string {
	if len(r.parentPath) == 0 {
		return fmt.Sprintf("%s%s.%s.fields-moved", r.blockKind.IDPrefix(), r.resourceType, r.childBlock)
	}
	return fmt.Sprintf("%s%s.%s.%s.fields-moved", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.parentPath, "."), r.childBlock)
}

func (r fieldMoveToNestedBlockRule) Description() string {
	if len(r.parentPath) == 0 {
		return fmt.Sprintf("Move %s fields into the %s block in v5", r.resourceType, r.childBlock)
	}
	return fmt.Sprintf("Move %s.%s fields into the %s block in v5", r.resourceType, strings.Join(r.parentPath, "."), r.childBlock)
}

func (r fieldMoveToNestedBlockRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		addr := r.blockKind.AddrPrefix() + rules.ResourceAddress(b)
		for _, parent := range rules.WalkNestedBlocks(b.Body(), r.parentPath) {
			var present []string
			for _, a := range r.attrs {
				if parent.GetAttribute(a) != nil {
					present = append(present, a)
				}
			}
			if len(present) > 0 {
				child := rules.FindOrCreateBlock(parent, r.childBlock)
				var moved []string
				for _, a := range present {
					if child.Body().GetAttribute(a) != nil {
						out = append(out, rules.Finding{
							RuleID:   r.ID(),
							Severity: rules.SeverityNeedsReview,
							Resource: addr,
							File:     ctx.File,
							Message:  fmt.Sprintf("cannot move %s into %s: that block already sets %s; reconcile manually", a, r.childBlock, a),
						})
						continue
					}
					rules.MoveAttribute(parent, child.Body(), a)
					moved = append(moved, a)
				}
				if len(moved) > 0 {
					out = append(out, rules.Finding{
						RuleID:   r.ID(),
						Severity: rules.SeverityChanged,
						Resource: addr,
						File:     ctx.File,
						Message:  fmt.Sprintf("moved %s into %s block", strings.Join(moved, ", "), r.childBlock),
					})
				}
			}

			child := parent.FirstMatchingBlock(r.childBlock, nil)
			if child == nil {
				continue
			}
			for _, rc := range r.requiredChildAttrs {
				if child.Body().GetAttribute(rc.name) != nil {
					continue
				}
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityNeedsReview,
					Resource: addr,
					File:     ctx.File,
					Message:  fmt.Sprintf("%s.%s is now required in v5: %s", r.childBlock, rc.name, rc.reason),
				})
			}

			if len(r.requiredChildAttrsWhenMultiple) > 0 {
				count := 0
				for _, blk := range parent.Blocks() {
					if blk.Type() == r.childBlock {
						count++
					}
				}
				if count > 1 {
					parentDesc := r.resourceType
					if len(r.parentPath) > 0 {
						parentDesc = r.parentPath[len(r.parentPath)-1]
					}
					for _, rc := range r.requiredChildAttrsWhenMultiple {
						for _, blk := range parent.Blocks() {
							if blk.Type() != r.childBlock || blk.Body().GetAttribute(rc.name) != nil {
								continue
							}
							out = append(out, rules.Finding{
								RuleID:   r.ID(),
								Severity: rules.SeverityNeedsReview,
								Resource: addr,
								File:     ctx.File,
								Message:  fmt.Sprintf("%s.%s is required in v5 when a %s has multiple %s blocks: %s", r.childBlock, rc.name, parentDesc, r.childBlock, rc.reason),
							})
						}
					}
				}
			}
		}
	}
	return out
}

// literalStringValue returns the resolvable literal string value of an
// hclwrite attribute's expression. ok is false when the expression isn't a
// plain template/literal string constant - e.g. it references a variable,
// local, or resource attribute, or calls a function - since those can't be
// evaluated without the caller's evaluation context.
func literalStringValue(attr *hclwrite.Attribute) (value string, ok bool) {
	src := attr.Expr().BuildTokens(nil).Bytes()
	expr, diags := hclsyntax.ParseExpression(src, "", hcl.InitialPos)
	if diags.HasErrors() {
		return "", false
	}
	switch expr.(type) {
	case *hclsyntax.TemplateExpr, *hclsyntax.LiteralValueExpr:
		// plain literal shapes only; fall through to evaluate below.
	default:
		return "", false
	}
	if len(expr.Variables()) > 0 {
		return "", false
	}
	v, diags := expr.Value(nil)
	if diags.HasErrors() || v.IsNull() || v.Type() != cty.String {
		return "", false
	}
	return v.AsString(), true
}

// literalBoolValue returns the literal boolean value of attr's expression.
// ok is false when the expression is not a bare true/false literal - e.g. it
// references a variable, local, or resource attribute, or calls a function -
// since those can't be evaluated without the caller's evaluation context.
func literalBoolValue(attr *hclwrite.Attribute) (value bool, ok bool) {
	src := attr.Expr().BuildTokens(nil).Bytes()
	expr, diags := hclsyntax.ParseExpression(src, "", hcl.InitialPos)
	if diags.HasErrors() {
		return false, false
	}
	switch expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		// bare true/false shape only; fall through to evaluate below.
	default:
		return false, false
	}
	if len(expr.Variables()) > 0 {
		return false, false
	}
	v, vdiags := expr.Value(nil)
	if vdiags.HasErrors() || v.IsNull() || v.Type() != cty.Bool {
		return false, false
	}
	return v.True(), true
}

// boolToEnumRule renames oldAttr → newAttr on every matched block (or nested
// block reachable by blockPath) where v5 replaces a boolean attribute with an
// enum-string attribute. When the value is a literal true/false, it is rewritten
// to trueValue/falseValue respectively (SeverityChanged). When the value is not
// a literal bool (var/local reference, interpolation, function call), the
// attribute is left byte-identical and flagged NeedsReview - there's no safe
// mechanical mapping, mirroring nameToIDRule's handling of non-traversal values.
//
// Note: rewriting the value via SetAttributeValue replaces the expression
// tokens, so an inline trailing comment on the old value's line is not
// preserved; lead comments are retained. Acceptable since the value changes.
type boolToEnumRule struct {
	blockKind    rules.BlockKind
	resourceType string
	blockPath    []string
	oldAttr      string
	newAttr      string
	trueValue    string
	falseValue   string
}

func (r boolToEnumRule) ID() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("%s%s.%s.bool-to-enum", r.blockKind.IDPrefix(), r.resourceType, r.oldAttr)
	}
	return fmt.Sprintf("%s%s.%s.%s.bool-to-enum", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.blockPath, "."), r.oldAttr)
}

func (r boolToEnumRule) Description() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("Migrate %s.%s (bool) to %s (enum) in v5", r.resourceType, r.oldAttr, r.newAttr)
	}
	return fmt.Sprintf("Migrate %s.%s.%s (bool) to %s (enum) in v5", r.resourceType, strings.Join(r.blockPath, "."), r.oldAttr, r.newAttr)
}

func (r boolToEnumRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		addr := r.blockKind.AddrPrefix() + rules.ResourceAddress(b)
		for _, nb := range rules.WalkNestedBlocks(b.Body(), r.blockPath) {
			attr := nb.GetAttribute(r.oldAttr)
			if attr == nil {
				continue
			}
			v, ok := literalBoolValue(attr)
			if !ok {
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityNeedsReview,
					Resource: addr,
					File:     ctx.File,
					Message:  fmt.Sprintf("%s was renamed to %s in v5 and its value changes from a bool to an enum (%q/%q); the current value is not a literal true/false, so update it manually", r.oldAttr, r.newAttr, r.trueValue, r.falseValue),
				})
				continue
			}
			enum := r.falseValue
			if v {
				enum = r.trueValue
			}
			rules.RenameAttribute(nb, r.oldAttr, r.newAttr)
			nb.SetAttributeValue(r.newAttr, cty.StringVal(enum))
			out = append(out, rules.Finding{
				RuleID:   r.ID(),
				Severity: rules.SeverityChanged,
				Resource: addr,
				File:     ctx.File,
				Message:  fmt.Sprintf("renamed %s to %s and mapped %t to %q", r.oldAttr, r.newAttr, v, enum),
			})
		}
	}
	return out
}

// literalStringList returns the resolvable literal []string value of an
// hclwrite attribute's expression. ok is false when the expression has any
// variable reference, isn't a list/tuple/set, or contains a non-string /
// non-literal element - since those can't be enumerated without the caller's
// evaluation context.
func literalStringList(attr *hclwrite.Attribute) ([]string, bool) {
	src := attr.Expr().BuildTokens(nil).Bytes()
	expr, diags := hclsyntax.ParseExpression(src, "", hcl.InitialPos)
	if diags.HasErrors() || len(expr.Variables()) > 0 {
		return nil, false
	}
	v, vdiags := expr.Value(nil)
	if vdiags.HasErrors() || v.IsNull() {
		return nil, false
	}
	if t := v.Type(); !t.IsListType() && !t.IsTupleType() && !t.IsSetType() {
		return nil, false
	}
	out := []string{}
	for it := v.ElementIterator(); it.Next(); {
		_, ev := it.Element()
		if ev.IsNull() || ev.Type() != cty.String {
			return nil, false
		}
		out = append(out, ev.AsString())
	}
	return out, true
}

// attrToBlockRule converts a list-of-strings attribute into repeated single-
// field blocks: `attr = ["a","b"]` becomes `newBlock { valueAttr = "a" }`
// `newBlock { valueAttr = "b" }`. When the value is not a literal string list
// (var/local reference, interpolation, or a list with a non-literal element),
// the attribute is left untouched and flagged NeedsReview - the elements can't
// be enumerated. New blocks are created with AppendNewBlock (multi-line, so no
// single-line gluing hazard) and appended after the parent's existing content.
type attrToBlockRule struct {
	blockKind    rules.BlockKind
	resourceType string
	blockPath    []string
	attr         string
	newBlock     string
	valueAttr    string
}

func (r attrToBlockRule) ID() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("%s%s.%s.attr-to-block", r.blockKind.IDPrefix(), r.resourceType, r.attr)
	}
	return fmt.Sprintf("%s%s.%s.%s.attr-to-block", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.blockPath, "."), r.attr)
}

func (r attrToBlockRule) Description() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("Convert %s.%s (list) to %s blocks in v5", r.resourceType, r.attr, r.newBlock)
	}
	return fmt.Sprintf("Convert %s.%s.%s (list) to %s blocks in v5", r.resourceType, strings.Join(r.blockPath, "."), r.attr, r.newBlock)
}

func (r attrToBlockRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		addr := r.blockKind.AddrPrefix() + rules.ResourceAddress(b)
		for _, parent := range rules.WalkNestedBlocks(b.Body(), r.blockPath) {
			attr := parent.GetAttribute(r.attr)
			if attr == nil {
				continue
			}
			elems, ok := literalStringList(attr)
			if !ok {
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityNeedsReview,
					Resource: addr,
					File:     ctx.File,
					Message:  fmt.Sprintf("%s became repeated %s blocks in v5; its value is not a literal list, so convert it to %s { %s = ... } blocks manually", r.attr, r.newBlock, r.newBlock, r.valueAttr),
				})
				continue
			}
			parent.RemoveAttribute(r.attr)
			for _, e := range elems {
				nb := parent.AppendNewBlock(r.newBlock, nil)
				nb.Body().SetAttributeValue(r.valueAttr, cty.StringVal(e))
			}
			out = append(out, rules.Finding{
				RuleID:   r.ID(),
				Severity: rules.SeverityChanged,
				Resource: addr,
				File:     ctx.File,
				Message:  fmt.Sprintf("converted %s to %d %s block(s)", r.attr, len(elems), r.newBlock),
			})
		}
	}
	return out
}

// negateFoldRule folds a cdn_frontdoor_rule condition's negate_condition into
// its operator, per the provider's own rule new_operator = (negate ? "Not" : "")
// + operator. When both operator (literal string) and negate_condition (literal
// bool) are present it rewrites operator and removes negate_condition (Changed);
// a literal-false negate_condition is simply removed (it was a no-op). When the
// values aren't literal (or negate is true with no literal operator), it leaves
// the block untouched and flags NeedsReview - the human must fold it manually.
//
// Note: rewriting operator via SetAttributeValue does not preserve an inline
// trailing comment on the operator line (parity with boolToEnumRule's note).
type negateFoldRule struct {
	blockKind    rules.BlockKind
	resourceType string
	blockPath    []string
	// skipOperatorValues names literal operator values that v5 itself removes
	// (e.g. "Any" on remote_address/socket_address). Folding negate_condition
	// into one of these would produce an operator value that's invalid on its
	// own (e.g. "NotAny") AND would remove negate_condition, hiding the
	// original value from the sibling valueChangeRule that flags it as
	// removed. When the literal operator is in this set, the fold is skipped
	// entirely (nothing edited) and a NeedsReview finding is emitted instead,
	// leaving both operator and negate_condition in place so the value-removed
	// flag still fires.
	skipOperatorValues []string
}

func (r negateFoldRule) ID() string {
	return fmt.Sprintf("%s%s.%s.negate_condition.negate-fold", r.blockKind.IDPrefix(), r.resourceType, strings.Join(r.blockPath, "."))
}

func (r negateFoldRule) Description() string {
	return fmt.Sprintf("Fold %s.%s.negate_condition into operator in v5", r.resourceType, strings.Join(r.blockPath, "."))
}

func (r negateFoldRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindBlocks(ctx.Write, r.blockKind, r.resourceType) {
		addr := r.blockKind.AddrPrefix() + rules.ResourceAddress(b)
		for _, nb := range rules.WalkNestedBlocks(b.Body(), r.blockPath) {
			neg := nb.GetAttribute("negate_condition")
			if neg == nil {
				continue
			}
			negVal, ok := literalBoolValue(neg)
			if !ok {
				out = append(out, rules.Finding{RuleID: r.ID(), Severity: rules.SeverityNeedsReview, Resource: addr, File: ctx.File,
					Message: "negate_condition is removed in v5 and its value is not a literal true/false; fold the negation into the operator (e.g. use a Not-prefixed operator) manually"})
				continue
			}
			if !negVal {
				nb.RemoveAttribute("negate_condition")
				out = append(out, rules.Finding{RuleID: r.ID(), Severity: rules.SeverityChanged, Resource: addr, File: ctx.File,
					Message: "removed negate_condition = false (no-op in v5)"})
				continue
			}
			op := nb.GetAttribute("operator")
			var opVal string
			okOp := false
			if op != nil {
				opVal, okOp = literalStringValue(op)
			}
			if !okOp {
				out = append(out, rules.Finding{RuleID: r.ID(), Severity: rules.SeverityNeedsReview, Resource: addr, File: ctx.File,
					Message: "negate_condition = true but operator is absent or not a literal string; set operator to the negated form (e.g. NotEqual) and remove negate_condition manually"})
				continue
			}
			if containsString(r.skipOperatorValues, opVal) {
				out = append(out, rules.Finding{RuleID: r.ID(), Severity: rules.SeverityNeedsReview, Resource: addr, File: ctx.File,
					Message: fmt.Sprintf("negate_condition = true but operator %q is removed in v5; choose a valid operator and fold the negation into it manually", opVal)})
				continue
			}
			nb.SetAttributeValue("operator", cty.StringVal("Not"+opVal))
			nb.RemoveAttribute("negate_condition")
			out = append(out, rules.Finding{RuleID: r.ID(), Severity: rules.SeverityChanged, Resource: addr, File: ctx.File,
				Message: fmt.Sprintf("folded negate_condition = true into operator (%q -> %q) and removed negate_condition", opVal, "Not"+opVal)})
		}
	}
	return out
}

func containsString(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}
