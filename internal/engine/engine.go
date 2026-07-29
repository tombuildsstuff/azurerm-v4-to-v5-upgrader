package engine

import (
	"bytes"
	"os"
	"strings"

	"github.com/pmezard/go-difflib/difflib"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/config"
	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/report"
	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

type Options struct {
	Dir       string
	DryRun    bool
	Recursive bool
	Rules     []rules.Rule
	// PreserveFormatting reverts hclwrite's incidental whitespace-only
	// re-formatting on lines the transforms did not semantically change, so the
	// diff contains only real changes. Off by default: normal output is
	// canonically formatted, matching `terraform fmt`.
	PreserveFormatting bool
}

type Result struct {
	Report *report.Report
	Diffs  map[string]string
}

func Run(opts Options) (*Result, error) {
	var files []*config.File
	var err error
	if opts.Recursive {
		files, err = config.LoadTree(opts.Dir)
	} else {
		files, err = config.Load(opts.Dir)
	}
	if err != nil {
		return nil, err
	}

	ruleset := opts.Rules
	if ruleset == nil {
		ruleset = rules.All()
	}

	rep := report.New()
	fired := map[string]bool{}
	diffs := map[string]string{}

	for _, f := range files {
		rep.ResourcesSeen += countAzureRMResources(f)
		ctx := &rules.RuleContext{File: f.Path, Write: f.Write.Body(), Syntax: f.Syntax}
		for _, r := range ruleset {
			for _, finding := range r.Apply(ctx) {
				if finding.File == "" {
					finding.File = f.Path
				}
				rep.Add(finding)
				fired[finding.RuleID] = true
			}
		}
	}
	rep.RulesFired = len(fired)

	for _, f := range files {
		out := f.Write.Bytes()
		if opts.PreserveFormatting {
			out = minimizeFormatChurn(f.Original, out, f.Path)
		}
		if bytes.Equal(out, f.Original) {
			continue
		}
		if opts.DryRun {
			diffs[f.Path] = unifiedDiff(f.Path, f.Original, out)
			continue
		}
		if err := os.WriteFile(f.Path, out, 0o644); err != nil {
			return nil, err
		}
	}

	return &Result{Report: rep, Diffs: diffs}, nil
}

func countAzureRMResources(f *config.File) int {
	n := 0
	for _, b := range f.Syntax.Blocks {
		if (b.Type == "resource" || b.Type == "data") && len(b.Labels) > 0 && strings.HasPrefix(b.Labels[0], "azurerm_") {
			n++
		}
	}
	return n
}

func unifiedDiff(path string, oldB, newB []byte) string {
	text, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(oldB)),
		B:        difflib.SplitLines(string(newB)),
		FromFile: path + " (v4)",
		ToFile:   path + " (v5)",
		Context:  3,
	})
	return text
}
