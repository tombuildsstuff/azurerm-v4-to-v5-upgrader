package azurerm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules/rulestest"
)

// findFixtureCaseDirs walks testdata and returns every directory that
// directly contains an input.tf fixture file.
func findFixtureCaseDirs(t *testing.T) []string {
	t.Helper()
	var dirs []string
	err := filepath.WalkDir("testdata", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "input.tf" {
			dirs = append(dirs, filepath.Dir(path))
		}
		return nil
	})
	require.NoError(t, err, "walking testdata")
	sort.Strings(dirs)
	return dirs
}

// TestFixtures runs every fixture case under testdata/** through the full
// default ruleset (rules.All()), comparing to that case's expected.tf and
// expected-report.json.
func TestFixtures(t *testing.T) {
	dirs := findFixtureCaseDirs(t)
	require.NotEmpty(t, dirs, "no fixture cases found under testdata")
	for _, dir := range dirs {
		dir := dir
		name, err := filepath.Rel("testdata", dir)
		require.NoError(t, err)
		t.Run(name, func(t *testing.T) {
			rulestest.RunGoldenAllRules(t, dir)
		})
	}
}

// TestEveryRuleIsExercised enforces that every rule registered in the global
// registry is fired by at least one fixture's expected-report.json. This
// replaces the old one-dir-per-rule requirement: a fixture may exercise many
// rules at once, and rule dirs need not be named after the rule ID.
func TestEveryRuleIsExercised(t *testing.T) {
	all := rules.All()
	require.NotEmpty(t, all)

	registered := make(map[string]struct{}, len(all))
	for _, r := range all {
		registered[r.ID()] = struct{}{}
	}

	exercised := map[string]struct{}{}
	dirs := findFixtureCaseDirs(t)
	require.NotEmpty(t, dirs, "no fixture cases found under testdata")
	for _, dir := range dirs {
		reportPath := filepath.Join(dir, "expected-report.json")
		b, err := os.ReadFile(reportPath)
		require.NoErrorf(t, err, "missing expected-report.json for fixture %s (run with -update)", dir)

		var report struct {
			Findings []struct {
				RuleID string `json:"rule_id"`
			} `json:"findings"`
		}
		require.NoErrorf(t, json.Unmarshal(b, &report), "parsing %s", reportPath)
		for _, f := range report.Findings {
			exercised[f.RuleID] = struct{}{}
		}
	}

	var missing []string
	for id := range registered {
		if _, ok := exercised[id]; !ok {
			missing = append(missing, id)
		}
	}
	sort.Strings(missing)
	require.Emptyf(t, missing, "rules registered but never exercised by any fixture: %v", missing)
}

// TestNoDuplicateRuleIDs guards the registry against two rules accidentally
// sharing the same ID, which would silently hide one of them from reports
// and from the coverage check above.
func TestNoDuplicateRuleIDs(t *testing.T) {
	all := rules.All()
	require.NotEmpty(t, all)

	seen := map[string][]int{}
	for i, r := range all {
		seen[r.ID()] = append(seen[r.ID()], i)
	}

	var collisions []string
	for id, idxs := range seen {
		if len(idxs) > 1 {
			collisions = append(collisions, fmt.Sprintf("%s (registered %d times)", id, len(idxs)))
		}
	}
	sort.Strings(collisions)
	require.Emptyf(t, collisions, "duplicate rule IDs in registry: %v", collisions)
}
