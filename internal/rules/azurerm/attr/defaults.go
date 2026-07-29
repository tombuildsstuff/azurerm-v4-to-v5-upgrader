package attr

import (
	"fmt"
	"strings"

	"github.com/zclconf/go-cty/cty"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

// changedDefaultRule pins a v5-changed default when the attribute is absent
// and the v4 default is known. When oldDefault is nil, the v4 default isn't
// known (e.g. it's computed, or simply hasn't been recorded yet); in that
// case the rule never guesses a value to inject, and instead surfaces a
// NeedsReview finding so a human can set the attribute explicitly.
type changedDefaultRule struct {
	resourceType string
	blockPath    []string
	attr         string
	oldDefault   *cty.Value
	newDefault   string
}

func (r changedDefaultRule) ID() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("%s.%s.default-changed", r.resourceType, r.attr)
	}
	return fmt.Sprintf("%s.%s.%s.default-changed", r.resourceType, strings.Join(r.blockPath, "."), r.attr)
}
func (r changedDefaultRule) Description() string {
	if len(r.blockPath) == 0 {
		return fmt.Sprintf("Pin v4 default for %s.%s (v5 changes it to %s)", r.resourceType, r.attr, r.newDefault)
	}
	return fmt.Sprintf("Pin v4 default for %s.%s.%s (v5 changes it to %s)", r.resourceType, strings.Join(r.blockPath, "."), r.attr, r.newDefault)
}
func (r changedDefaultRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	var out []rules.Finding
	for _, b := range rules.FindResourceBlocks(ctx.Write, r.resourceType) {
		for _, nb := range rules.WalkNestedBlocks(b.Body(), r.blockPath) {
			if nb.GetAttribute(r.attr) != nil {
				continue
			}
			if r.oldDefault == nil {
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityNeedsReview,
					Resource: rules.ResourceAddress(b),
					File:     ctx.File,
					Message:  fmt.Sprintf("v5 changes the default of %s to '%s'; set it explicitly to preserve v4 behavior", r.attr, r.newDefault),
				})
				continue
			}
			if rules.PinChangedDefault(nb, r.attr, *r.oldDefault, r.newDefault) {
				out = append(out, rules.Finding{
					RuleID:   r.ID(),
					Severity: rules.SeverityNote,
					Resource: rules.ResourceAddress(b),
					File:     ctx.File,
					Message:  fmt.Sprintf("pinned %s to v4 default; v5 default is now %s", r.attr, r.newDefault),
				})
			}
		}
	}
	return out
}
