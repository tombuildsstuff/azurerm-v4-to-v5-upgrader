package githubactions

import (
	"fmt"
	"io"
	"strings"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

// WriteAnnotations prints one GitHub Actions workflow-command annotation per
// finding to w. NeedsReview findings become ::error, all others ::warning.
// The line property is omitted when the finding has no line.
func WriteAnnotations(w io.Writer, findings []rules.Finding) error {
	for _, f := range findings {
		level := "warning"
		if f.Severity == rules.SeverityNeedsReview {
			level = "error"
		}
		var err error
		if f.Range.Start.Line > 0 {
			_, err = fmt.Fprintf(w, "::%s file=%s,line=%d::%s\n",
				level, escapeProperty(f.File), f.Range.Start.Line, escapeData(f.Message))
		} else {
			_, err = fmt.Fprintf(w, "::%s file=%s::%s\n",
				level, escapeProperty(f.File), escapeData(f.Message))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// escapeData escapes a workflow-command message body per GitHub's rules.
func escapeData(s string) string {
	r := strings.NewReplacer("%", "%25", "\r", "%0D", "\n", "%0A")
	return r.Replace(s)
}

// escapeProperty escapes a workflow-command property value (e.g. file) per
// GitHub's rules - a superset of escapeData that also escapes ',' and ':'.
func escapeProperty(s string) string {
	r := strings.NewReplacer("%", "%25", "\r", "%0D", "\n", "%0A", ":", "%3A", ",", "%2C")
	return r.Replace(s)
}
