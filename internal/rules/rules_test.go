package rules

import "testing"

func TestSeverityString(t *testing.T) {
	cases := map[Severity]string{
		SeverityChanged:     "Changed",
		SeverityNote:        "Note",
		SeverityNeedsReview: "NeedsReview",
	}
	for sev, want := range cases {
		if got := sev.String(); got != want {
			t.Errorf("Severity(%d).String() = %q, want %q", sev, got, want)
		}
	}
}

type fakeRule struct{ id string }

func (f fakeRule) ID() string                   { return f.id }
func (f fakeRule) Description() string          { return "fake" }
func (f fakeRule) Apply(*RuleContext) []Finding { return nil }

func TestRegistryPreservesOrder(t *testing.T) {
	resetRegistry()
	Register(fakeRule{id: "a"})
	Register(fakeRule{id: "b"})
	all := All()
	if len(all) != 2 || all[0].ID() != "a" || all[1].ID() != "b" {
		t.Fatalf("unexpected registry contents: %+v", all)
	}
}
