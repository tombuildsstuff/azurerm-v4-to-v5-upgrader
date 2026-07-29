package attr

import (
	"strings"
	"testing"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules/rulestest"
)

// TestRenamedAttrRuleDataSource proves the shared kinds target `data` blocks
// when blockKind is DataBlock, and that the ID gains the `data.` prefix.
func TestRenamedAttrRuleDataSource(t *testing.T) {
	r := renamedAttrRule{
		blockKind:    rules.DataBlock,
		resourceType: "azurerm_key_vault",
		oldAttr:      "enable_rbac_authorization",
		newAttr:      "rbac_authorization_enabled",
	}
	rulestest.RunGolden(t, r, "testdata/data-source-rename")
}

// TestRenamedBlockRuleTopLevel covers: basic rename + inner attr renames,
// repeated sibling blocks, an absent old block (no-op), and idempotency.
func TestRenamedBlockRuleTopLevel(t *testing.T) {
	r := renamedBlockRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		oldBlock:     "request_header_action",
		newBlock:     "modify_request_header",
		innerRenames: map[string]string{"header_action": "operator", "value": "header_value"},
	}
	rulestest.RunGolden(t, r, "testdata/renamed-block/top-level")
}

// TestRenamedBlockRuleNested covers a rename scoped under a parent block path,
// plus an inner attr removal.
func TestRenamedBlockRuleNested(t *testing.T) {
	r := renamedBlockRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		parentPath:   []string{"conditions"},
		oldBlock:     "cookies_condition",
		newBlock:     "request_cookies",
		innerRenames: map[string]string{"match_values": "values"},
		innerRemoves: []string{"negate_condition"},
	}
	rulestest.RunGolden(t, r, "testdata/renamed-block/nested")
}

// TestFieldMoveToNestedBlockRule covers: child created when absent (+ moved
// attrs + now-required name flagged), child reused when present (name present,
// no flag), a comment-carrying attr, a subset of attrs present, and
// idempotency.
func TestFieldMoveToNestedBlockRule(t *testing.T) {
	r := fieldMoveToNestedBlockRule{
		resourceType: "azurerm_site_recovery_replicated_vm",
		parentPath:   []string{"network_interface"},
		childBlock:   "ip_configuration",
		attrs:        []string{"target_static_ip", "target_subnet_name", "failover_test_static_ip"},
		requiredChildAttrs: []requiredChildAttr{
			{name: "name", reason: "must match the source VM's NIC IP configuration name"},
		},
	}
	rulestest.RunGolden(t, r, "testdata/field-move")
}

// TestFieldMoveToNestedBlockRuleTopLevelID proves ID() guards an empty
// parentPath: without the guard, strings.Join(nil, ".") yields "" and the
// unconditional format produces a malformed double-dot ID like
// "azurerm_x..ip_configuration.fields-moved".
func TestFieldMoveToNestedBlockRuleTopLevelID(t *testing.T) {
	r := fieldMoveToNestedBlockRule{
		resourceType: "azurerm_x",
		childBlock:   "ip_configuration",
		attrs:        []string{"foo"},
	}
	id := r.ID()
	if id != "azurerm_x.ip_configuration.fields-moved" {
		t.Errorf("ID() = %q, want %q", id, "azurerm_x.ip_configuration.fields-moved")
	}
	if strings.Contains(id, "..") {
		t.Errorf("ID() = %q contains a double dot", id)
	}
}
