package githubactions

import (
	"fmt"
	"io"
)

// WriteOutputs appends the action outputs (needs-changes, findings-count) to w,
// which is normally the file at $GITHUB_OUTPUT.
func WriteOutputs(w io.Writer, needsChanges bool, findingsCount int) error {
	_, err := fmt.Fprintf(w, "needs-changes=%t\nfindings-count=%d\n", needsChanges, findingsCount)
	return err
}
