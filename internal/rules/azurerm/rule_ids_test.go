package azurerm

import (
	"flag"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

var updateRuleIDs = flag.Bool("update-rule-ids", false, "rewrite the rule-ids golden file")

// TestRegisteredRuleIDsSnapshot pins the exact set of registered rule IDs.
// Adding a BlockKind field (Task 4) must not shift any existing ID: since
// ResourceBlock is the zero value and its IDPrefix() is "", this snapshot must
// stay byte-identical through the refactor. When new rules are intentionally
// registered later, regenerate with -update-rule-ids and review the diff.
func TestRegisteredRuleIDsSnapshot(t *testing.T) {
	all := rules.All()
	require.NotEmpty(t, all)

	ids := make([]string, 0, len(all))
	for _, r := range all {
		ids = append(ids, r.ID())
	}
	sort.Strings(ids)
	got := strings.Join(ids, "\n") + "\n"

	const golden = "testdata/rule-ids.golden"
	if *updateRuleIDs {
		require.NoError(t, os.WriteFile(golden, []byte(got), 0o644))
		return
	}
	want, err := os.ReadFile(golden)
	require.NoError(t, err, "missing testdata/rule-ids.golden (run with -update-rule-ids)")
	require.Equal(t, string(want), got, "registered rule IDs changed unexpectedly")
}
