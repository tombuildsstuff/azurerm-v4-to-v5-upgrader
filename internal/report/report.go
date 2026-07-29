package report

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

// Report aggregates findings and coverage counters from one engine run.
type Report struct {
	Findings      []rules.Finding
	ResourcesSeen int
	RulesFired    int
}

func New() *Report { return &Report{} }

func (r *Report) Add(f rules.Finding) { r.Findings = append(r.Findings, f) }

// ExitCode returns 1 when the run needs human attention, else 0.
func (r *Report) ExitCode() int {
	for _, f := range r.Findings {
		if f.Severity == rules.SeverityNeedsReview {
			return 1
		}
	}
	return 0
}

func (r *Report) sorted() []rules.Finding {
	out := append([]rules.Finding(nil), r.Findings...)
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Severity != b.Severity {
			return a.Severity < b.Severity
		}
		if a.RuleID != b.RuleID {
			return a.RuleID < b.RuleID
		}
		if a.File != b.File {
			return a.File < b.File
		}
		return a.Range.Start.Line < b.Range.Start.Line
	})
	return out
}

type jsonFinding struct {
	RuleID   string `json:"rule_id"`
	Severity string `json:"severity"`
	Resource string `json:"resource"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Message  string `json:"message"`
}

type jsonReport struct {
	Findings      []jsonFinding `json:"findings"`
	ResourcesSeen int           `json:"resources_seen"`
	RulesFired    int           `json:"rules_fired"`
}

func (r *Report) RenderJSON(w io.Writer) error {
	out := jsonReport{
		Findings:      []jsonFinding{},
		ResourcesSeen: r.ResourcesSeen,
		RulesFired:    r.RulesFired,
	}
	for _, f := range r.sorted() {
		out.Findings = append(out.Findings, jsonFinding{
			RuleID:   f.RuleID,
			Severity: f.Severity.String(),
			Resource: f.Resource,
			File:     f.File,
			Line:     f.Range.Start.Line,
			Message:  f.Message,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// jsonReportValidate mirrors jsonReport but adds a top-level needs_upgrade
// field. It is used solely by RenderJSONValidate so that update's RenderJSON
// output shape is left untouched.
type jsonReportValidate struct {
	jsonReport
	NeedsUpgrade bool `json:"needs_upgrade"`
}

// RenderJSONValidate renders the report as JSON with an added top-level
// "needs_upgrade" field, for use by the read-only validate subcommand.
func (r *Report) RenderJSONValidate(w io.Writer, needsUpgrade bool) error {
	out := jsonReportValidate{
		jsonReport: jsonReport{
			Findings:      []jsonFinding{},
			ResourcesSeen: r.ResourcesSeen,
			RulesFired:    r.RulesFired,
		},
		NeedsUpgrade: needsUpgrade,
	}
	for _, f := range r.sorted() {
		out.Findings = append(out.Findings, jsonFinding{
			RuleID:   f.RuleID,
			Severity: f.Severity.String(),
			Resource: f.Resource,
			File:     f.File,
			Line:     f.Range.Start.Line,
			Message:  f.Message,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func (r *Report) RenderText(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "azurerm v5 upgrade: %d resource(s) seen, %d rule(s) fired\n\n",
		r.ResourcesSeen, r.RulesFired); err != nil {
		return err
	}
	for _, f := range r.sorted() {
		loc := f.File
		if f.Range.Start.Line > 0 {
			loc = fmt.Sprintf("%s:%d", f.File, f.Range.Start.Line)
		}
		if _, err := fmt.Fprintf(w, "[%s] %s (%s) %s\n",
			f.Severity, f.RuleID, loc, f.Message); err != nil {
			return err
		}
	}
	return nil
}

// escapeMarkdownCell makes s safe to place in a Markdown table cell: pipes are
// escaped and newlines collapsed to spaces so a single row stays intact.
func escapeMarkdownCell(s string) string {
	r := strings.NewReplacer("|", "\\|", "\r\n", " ", "\n", " ", "\r", " ")
	return r.Replace(s)
}

// RenderMarkdown renders the report as GitHub-flavored markdown, used for both
// the Actions step summary and the PR comment. needsUpgrade drives the headline.
func (r *Report) RenderMarkdown(w io.Writer, needsUpgrade bool) error {
	if _, err := fmt.Fprintf(w, "## azurerm v5 upgrade check\n\n"); err != nil {
		return err
	}
	if !needsUpgrade {
		_, err := fmt.Fprintf(w, "Changes required: **no** - module is already on azurerm v5.\n")
		return err
	}
	if _, err := fmt.Fprintf(w, "Changes required: **yes** (%d resource(s) seen, %d rule(s) fired).\n\n",
		r.ResourcesSeen, r.RulesFired); err != nil {
		return err
	}
	groups := map[string][]rules.Finding{}
	for _, f := range r.sorted() {
		dir := filepath.Dir(f.File)
		groups[dir] = append(groups[dir], f)
	}
	dirs := make([]string, 0, len(groups))
	for dir := range groups {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	for _, dir := range dirs {
		if _, err := fmt.Fprintf(w, "### %s\n\n", dir); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "| Severity | Rule | Location | Message |\n| --- | --- | --- | --- |\n"); err != nil {
			return err
		}
		for _, f := range groups[dir] {
			loc := f.File
			if f.Range.Start.Line > 0 {
				loc = fmt.Sprintf("%s:%d", f.File, f.Range.Start.Line)
			}
			if _, err := fmt.Fprintf(w, "| %s | %s | %s | %s |\n",
				escapeMarkdownCell(f.Severity.String()), escapeMarkdownCell(f.RuleID), escapeMarkdownCell(loc), escapeMarkdownCell(f.Message)); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}
