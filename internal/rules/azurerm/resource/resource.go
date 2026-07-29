package resource

import (
	"fmt"
	"sort"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

// removedResourceRule flags resource types that no longer exist in v5.
// It never edits config - removed resources need manual migration.
type removedResourceRule struct {
	types map[string]struct{}
}

func (removedResourceRule) ID() string          { return "resource.removed" }
func (removedResourceRule) Description() string { return "Flag resource types removed in AzureRM v5" }

func (r removedResourceRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	// Deterministic order: collect sorted type list, scan each.
	typesSorted := make([]string, 0, len(r.types))
	for t := range r.types {
		typesSorted = append(typesSorted, t)
	}
	sort.Strings(typesSorted)

	var out []rules.Finding
	for _, t := range typesSorted {
		for _, b := range rules.FindResourceBlocks(ctx.Write, t) {
			out = append(out, rules.Finding{
				RuleID:   "resource.removed",
				Severity: rules.SeverityNeedsReview,
				Resource: rules.ResourceAddress(b),
				File:     ctx.File,
				Message:  fmt.Sprintf("%s is removed in AzureRM v5 and cannot be auto-migrated; replace it manually", t),
			})
		}
	}
	return out
}
