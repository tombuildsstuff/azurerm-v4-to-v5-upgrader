package resource

import (
	"fmt"
	"sort"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

// removedDataSourceRule flags `data "azurerm_..."` blocks whose type no longer
// exists in v5. Like removedResourceRule, it never edits config.
type removedDataSourceRule struct {
	types map[string]struct{}
}

func (removedDataSourceRule) ID() string { return "data-source.removed" }
func (removedDataSourceRule) Description() string {
	return "Flag data source types removed in AzureRM v5"
}

func (r removedDataSourceRule) Apply(ctx *rules.RuleContext) []rules.Finding {
	typesSorted := make([]string, 0, len(r.types))
	for t := range r.types {
		typesSorted = append(typesSorted, t)
	}
	sort.Strings(typesSorted)

	var out []rules.Finding
	for _, t := range typesSorted {
		for _, b := range rules.FindBlocks(ctx.Write, rules.DataBlock, t) {
			out = append(out, rules.Finding{
				RuleID:   "data-source.removed",
				Severity: rules.SeverityNeedsReview,
				Resource: rules.DataBlock.AddrPrefix() + rules.ResourceAddress(b),
				File:     ctx.File,
				Message:  fmt.Sprintf("the %s data source is removed in AzureRM v5 and cannot be auto-migrated; replace it manually", t),
			})
		}
	}
	return out
}
