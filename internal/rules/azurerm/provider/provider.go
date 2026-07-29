package provider

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

const pinVersion = "=5.0.0"

type versionPinRule struct{}

// NewVersionPinRule pins the azurerm provider to =5.0.0.
func NewVersionPinRule() rules.Rule { return versionPinRule{} }

func (versionPinRule) ID() string          { return "provider.version-pin" }
func (versionPinRule) Description() string { return "Pin the azurerm provider to =5.0.0" }

func (r versionPinRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, tfBlock := range ctx.Write.Blocks() {
		if tfBlock.Type() != "terraform" {
			continue
		}
		for _, rpBlock := range tfBlock.Body().Blocks() {
			if rpBlock.Type() != "required_providers" {
				continue
			}
			attr := rpBlock.Body().GetAttribute("azurerm")
			if attr == nil {
				continue
			}
			if setProviderVersion(rpBlock.Body(), attr) {
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityChanged,
					File:     ctx.File,
					Range:    attrRange(ctx, "terraform"),
					Message:  "pinned azurerm provider to " + pinVersion,
				})
			}
		}
	}
	return out
}

// setProviderVersion surgically rewrites only the value of the version key
// inside the azurerm = { ... } object expression. It never touches source,
// configuration_aliases, any other key, or comments - those are preserved
// byte-for-byte, since hclwrite has no structured object-item API and a
// naive SetAttributeValue rebuild of the whole object would silently drop
// everything except the keys it's told about (this previously deleted
// configuration_aliases, which reusable modules rely on for provider
// aliasing, plus any inline comments).
//
// It returns false (no change made) when:
//   - the azurerm attribute's value is not a recognized object constructor
//     (e.g. a traversal), so the rule must not touch a target it can't verify, or
//   - the version is already pinned to pinVersion, making the run idempotent.
func setProviderVersion(rpBody *hclwrite.Body, attr *hclwrite.Attribute) bool {
	oce, tokens, ok := parseObjectExpr(attr)
	if !ok {
		// Not a recognized object shape (e.g. a traversal) - don't touch it.
		return false
	}

	pinTokens := hclwrite.TokensForValue(cty.StringVal(pinVersion))

	if item, found := findObjectItem(oce, "version"); found {
		if vv, diags := item.ValueExpr.Value(nil); !diags.HasErrors() && vv.Type().FriendlyName() == "string" && vv.AsString() == pinVersion {
			// Already pinned - nothing to do.
			return false
		}
		newTokens := replaceExprRange(tokens, item.ValueExpr.Range(), pinTokens)
		rpBody.SetAttributeRaw("azurerm", newTokens)
		return true
	}

	// No version key present at all (unusual, but possible) - insert one
	// rather than touching anything that's already there.
	newTokens := insertObjectItem(tokens, oce, "version", pinTokens)
	rpBody.SetAttributeRaw("azurerm", newTokens)
	return true
}

func attrRange(ctx *rules.RuleContext, blockType string) hcl.Range {
	for _, b := range ctx.Syntax.Blocks {
		if b.Type == blockType {
			return b.DefRange()
		}
	}
	return hcl.Range{}
}
