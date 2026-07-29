package githubactions

import (
	"fmt"
	"io"
	"net/http"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

// Deps carries the injectable seams Emit needs. Summary and Output may be nil
// when their target env var is unset; they are then skipped.
type Deps struct {
	Env        Environment
	Stdout     io.Writer
	Summary    io.Writer
	Output     io.Writer
	ReadFile   func(string) ([]byte, error)
	HTTPClient *http.Client
}

// Emit renders every applicable GitHub Actions channel for a detection result.
// markdown is the pre-rendered body shared by the step summary and the PR
// comment. It returns an error only on operational failure: a missing token in
// a pull-request context, an unreadable/invalid event payload, or a failed API
// POST. Annotations are emitted only inside Actions; summary/output only when a
// target writer is present; a comment only in a pull-request context.
func Emit(d Deps, findings []rules.Finding, markdown string, needsChanges bool) error {
	if d.Env.IsActions {
		if err := WriteAnnotations(d.Stdout, findings); err != nil {
			return err
		}
	}
	if d.Summary != nil {
		if _, err := io.WriteString(d.Summary, markdown); err != nil {
			return err
		}
	}
	if d.Output != nil {
		if err := WriteOutputs(d.Output, needsChanges, len(findings)); err != nil {
			return err
		}
	}
	if d.Env.IsPullRequest() {
		if d.Env.Token == "" {
			return fmt.Errorf("GITHUB_TOKEN is required to comment on pull requests")
		}
		payload, err := d.ReadFile(d.Env.EventPath)
		if err != nil {
			return fmt.Errorf("reading event payload: %w", err)
		}
		prNumber, err := ParsePRNumber(payload)
		if err != nil {
			return err
		}
		if err := PostComment(d.HTTPClient, d.Env, prNumber, markdown); err != nil {
			return err
		}
	}
	return nil
}
