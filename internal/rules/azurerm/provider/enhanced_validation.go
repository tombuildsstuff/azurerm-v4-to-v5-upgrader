package provider

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

// enhancedValidationRule migrates the provider's enhanced_validation block:
// v5 moves it from a top-level provider attribute-block into the features
// block, and flips the defaults of its locations/resource_providers sub-fields
// from enabled (v4) to disabled (v5). This rule pins locations/resource_providers
// to true when absent (preserving v4 behavior), then relocates a top-level
// enhanced_validation block into features (creating features if absent). If the
// block already lives inside features, only the default-pinning applies.
//
// The relocation preserves the block's contents and comments by moving its full
// token stream; pinning happens before the move, while the block is still a
// structured node.
type enhancedValidationRule struct{}

// NewEnhancedValidationRule returns the rule.
func NewEnhancedValidationRule() rules.Rule { return enhancedValidationRule{} }

func (enhancedValidationRule) ID() string { return "provider.enhanced-validation" }
func (enhancedValidationRule) Description() string {
	return "Move enhanced_validation into the features block and pin its flipped defaults in v5"
}

func (r enhancedValidationRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, pb := range azureRMProviderBlocks(ctx.Write) {
		body := pb.Body()

		// Locate enhanced_validation: prefer a top-level block (the deprecated
		// v4 location that must move); otherwise look inside features.
		topEV := body.FirstMatchingBlock("enhanced_validation", nil)
		ev := topEV
		if ev == nil {
			if feat := body.FirstMatchingBlock("features", nil); feat != nil {
				ev = feat.Body().FirstMatchingBlock("enhanced_validation", nil)
			}
		}
		if ev == nil {
			continue
		}

		// Pin flipped defaults (only when absent) while ev is still structured.
		pinned := pinBoolDefault(ev.Body(), "locations", true)
		if pinBoolDefault(ev.Body(), "resource_providers", true) {
			pinned = true
		}

		moved := false
		if topEV != nil {
			feat := rules.FindOrCreateBlock(body, "features")
			// features may be single-line ("features {}"), so ensure the
			// moved block starts on its own line rather than gluing onto
			// the opening brace, which would otherwise produce invalid
			// HCL ("A single-line block definition can contain only a
			// single argument").
			feat.Body().AppendNewline()
			feat.Body().AppendUnstructuredTokens(ev.BuildTokens(nil))
			body.RemoveBlock(ev)
			moved = true
		}

		if moved || pinned {
			msg := "pinned enhanced_validation.locations/resource_providers = true to preserve v4 defaults"
			if moved {
				msg = "moved enhanced_validation into the features block and " + msg
			}
			out = append(out, rules.Finding{
				RuleID:   r.ID(),
				Severity: rules.SeverityChanged,
				File:     ctx.File,
				Message:  msg,
			})
		}
	}
	return out
}

// pinBoolDefault sets name = val on body when name is absent, returning true if
// it added the attribute. An existing (explicit) value is left untouched.
func pinBoolDefault(body *hclwrite.Body, name string, val bool) bool {
	if body.GetAttribute(name) != nil {
		return false
	}
	body.SetAttributeValue(name, cty.BoolVal(val))
	return true
}
