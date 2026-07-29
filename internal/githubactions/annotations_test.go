package githubactions

import (
	"bytes"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

func TestWriteAnnotationsErrorAndWarning(t *testing.T) {
	findings := []rules.Finding{
		{RuleID: "a", Severity: rules.SeverityNeedsReview, File: "main.tf",
			Range: hcl.Range{Start: hcl.Pos{Line: 5}}, Message: "must review"},
		{RuleID: "b", Severity: rules.SeverityChanged, File: "vars.tf",
			Range: hcl.Range{Start: hcl.Pos{Line: 0}}, Message: "changed it"},
	}
	var buf bytes.Buffer
	require.NoError(t, WriteAnnotations(&buf, findings))
	out := buf.String()
	require.Contains(t, out, "::error file=main.tf,line=5::must review\n")
	require.Contains(t, out, "::warning file=vars.tf::changed it\n")
}

func TestWriteAnnotationsEscapesNewlines(t *testing.T) {
	findings := []rules.Finding{
		{RuleID: "a", Severity: rules.SeverityNeedsReview, File: "main.tf",
			Range: hcl.Range{Start: hcl.Pos{Line: 1}}, Message: "line1\nline2"},
	}
	var buf bytes.Buffer
	require.NoError(t, WriteAnnotations(&buf, findings))
	require.Contains(t, buf.String(), "line1%0Aline2")
}

func TestWriteOutputs(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, WriteOutputs(&buf, true, 3))
	require.Equal(t, "needs-changes=true\nfindings-count=3\n", buf.String())
}
