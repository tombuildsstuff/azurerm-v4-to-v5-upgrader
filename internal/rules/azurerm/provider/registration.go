package provider

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

// providerRegistrationRule reconciles the v5 provider-registration changes on
// each `provider "azurerm"` block: skip_provider_registration was removed and
// resource_provider_registrations' default flipped from "legacy" to "none".
// To preserve v4 behavior:
//   - skip = true  -> resource_provider_registrations = "none"  (registered nothing)
//   - skip = false -> resource_provider_registrations = "legacy"
//   - skip absent + registrations absent -> pin "legacy" (the old default)
//
// skip is always removed. Ambiguous states (non-literal skip, or both skip and
// an explicit registrations value) are flagged rather than guessed.
type providerRegistrationRule struct{}

// NewProviderRegistrationRule returns the rule.
func NewProviderRegistrationRule() rules.Rule { return providerRegistrationRule{} }

func (providerRegistrationRule) ID() string { return "provider.resource-provider-registrations" }
func (providerRegistrationRule) Description() string {
	return "Reconcile skip_provider_registration and the resource_provider_registrations default change in v5"
}

func (r providerRegistrationRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, pb := range azureRMProviderBlocks(ctx.Write) {
		body := pb.Body()
		skip := body.GetAttribute("skip_provider_registration")
		reg := body.GetAttribute("resource_provider_registrations")

		switch {
		case skip != nil:
			b, ok := literalBool(skip)
			if !ok {
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityNeedsReview,
					File:     ctx.File,
					Message:  "skip_provider_registration was removed in v5; its value is not a literal true/false, so set resource_provider_registrations manually (\"none\" for true, \"legacy\" for false)",
				})
				continue
			}
			if reg != nil {
				body.RemoveAttribute("skip_provider_registration")
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityNeedsReview,
					File:     ctx.File,
					Message:  "removed deprecated skip_provider_registration; confirm the existing resource_provider_registrations value is what you want",
				})
				continue
			}
			val := "legacy"
			if b {
				val = "none"
			}
			body.RemoveAttribute("skip_provider_registration")
			body.SetAttributeValue("resource_provider_registrations", cty.StringVal(val))
			out = append(out, rules.Finding{
				RuleID:   r.ID(),
				Severity: rules.SeverityChanged,
				File:     ctx.File,
				Message:  fmt.Sprintf("translated skip_provider_registration = %t to resource_provider_registrations = %q", b, val),
			})
		case reg == nil:
			body.SetAttributeValue("resource_provider_registrations", cty.StringVal("legacy"))
			out = append(out, rules.Finding{
				RuleID:   r.ID(),
				Severity: rules.SeverityNote,
				File:     ctx.File,
				Message:  "pinned resource_provider_registrations = \"legacy\" (the v5 default changed to \"none\"; \"legacy\" preserves v4 behavior)",
			})
		default:
			// resource_provider_registrations present, skip absent: nothing to do.
		}
	}
	return out
}
