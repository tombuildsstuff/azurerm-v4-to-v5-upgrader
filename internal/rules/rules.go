package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Severity ranks a finding. Higher values are more urgent.
type Severity int

const (
	SeverityChanged     Severity = iota // an edit was applied; behavior preserved
	SeverityNote                        // informational edit (e.g. pinned changed default)
	SeverityNeedsReview                 // a human must decide; forces non-zero exit
)

func (s Severity) String() string {
	switch s {
	case SeverityChanged:
		return "Changed"
	case SeverityNote:
		return "Note"
	case SeverityNeedsReview:
		return "NeedsReview"
	default:
		return "Unknown"
	}
}

// Finding records one thing a rule observed or did.
type Finding struct {
	RuleID   string
	Severity Severity
	Resource string
	File     string
	Range    hcl.Range
	Message  string
}

// RuleContext gives a rule both views of a single file's top-level body.
type RuleContext struct {
	File   string
	Write  *hclwrite.Body
	Syntax *hclsyntax.Body
}

// Rule is one hand-authored v5 transformation.
type Rule interface {
	ID() string
	Description() string
	Apply(ctx *RuleContext) []Finding
}
