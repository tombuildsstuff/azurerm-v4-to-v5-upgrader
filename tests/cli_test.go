package tests

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "rewrite golden fixture files from the binary's output")

// binaryPath is the freshly built upgrader binary, shared by every test.
var binaryPath string

const corpusDir = "testdata/corpus"

func TestMain(m *testing.M) {
	flag.Parse()
	dir, err := os.MkdirTemp("", "upgrader-bin")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	binaryPath = filepath.Join(dir, "upgrader")
	build := exec.Command("go", "build", "-o", binaryPath, "github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		panic("building upgrader binary: " + err.Error())
	}
	os.Exit(m.Run())
}

// runUpgrader execs `upgrader update` recursively over dir and returns stdout
// (the JSON report) and the process exit code. A NeedsReview finding makes the
// binary exit 1, which is expected, so a non-zero exit is not an error here.
func runUpgrader(t *testing.T, dir string) (stdout string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binaryPath, "update", "--dir", dir, "--report", "json")
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return out.String(), ee.ExitCode()
		}
		t.Fatalf("running upgrader: %v\nstderr: %s", err, errBuf.String())
	}
	return out.String(), 0
}

// copyTree copies every file under src into dst, preserving structure.
func copyTree(t *testing.T, src, dst string) {
	t.Helper()
	require.NoError(t, filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, b, 0o644)
	}))
}

// tfFiles returns the sorted repo-relative paths of all .tf files under root.
func tfFiles(t *testing.T, root string) []string {
	t.Helper()
	var files []string
	require.NoError(t, filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".tf") {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files = append(files, rel)
		}
		return nil
	}))
	slices.Sort(files)
	return files
}

// normalizeReport strips the temp dir prefix so "file" fields are stable,
// repo-relative paths (e.g. modules/storage/main.tf).
func normalizeReport(s, dir string) string {
	return strings.ReplaceAll(s, dir+string(os.PathSeparator), "")
}

// severities extracts the set of finding severities from a report body.
func severities(t *testing.T, reportJSON string) []string {
	t.Helper()
	var r struct {
		Findings []struct {
			Severity string `json:"severity"`
		} `json:"findings"`
	}
	require.NoError(t, json.Unmarshal([]byte(reportJSON), &r))
	out := make([]string, 0, len(r.Findings))
	for _, f := range r.Findings {
		out = append(out, f.Severity)
	}
	return out
}

// assertExitCode checks the CLI exit-code contract: 1 iff a NeedsReview
// finding is present, else 0.
func assertExitCode(t *testing.T, exitCode int, sevs []string) {
	t.Helper()
	want := 0
	if slices.Contains(sevs, "NeedsReview") {
		want = 1
	}
	require.Equal(t, want, exitCode, "exit code should reflect NeedsReview presence")
}

func TestGolden(t *testing.T) {
	tmp := t.TempDir()
	copyTree(t, filepath.Join(corpusDir, "input"), tmp)

	reportOut, exitCode := runUpgrader(t, tmp)
	normReport := normalizeReport(reportOut, tmp)
	assertExitCode(t, exitCode, severities(t, normReport))

	// Every upgraded .tf must parse as valid HCL.
	for _, rel := range tfFiles(t, tmp) {
		b, err := os.ReadFile(filepath.Join(tmp, rel))
		require.NoError(t, err)
		_, diags := hclsyntax.ParseConfig(b, rel, hcl.InitialPos)
		require.False(t, diags.HasErrors(), "%s: invalid HCL: %s", rel, diags.Error())
	}

	expectedDir := filepath.Join(corpusDir, "expected")
	if *update {
		require.NoError(t, os.RemoveAll(expectedDir))
		copyTree(t, tmp, expectedDir)
		require.NoError(t, os.WriteFile(filepath.Join(corpusDir, "expected-report.json"), []byte(normReport), 0o644))
		return
	}

	// Whole-tree golden: identical file set and identical contents.
	wantFiles := tfFiles(t, expectedDir)
	gotFiles := tfFiles(t, tmp)
	require.Equal(t, wantFiles, gotFiles, "file set mismatch vs expected/")
	for _, rel := range wantFiles {
		want, err := os.ReadFile(filepath.Join(expectedDir, rel))
		require.NoError(t, err)
		got, err := os.ReadFile(filepath.Join(tmp, rel))
		require.NoError(t, err)
		require.Equal(t, string(want), string(got), "content mismatch: %s", rel)
	}

	wantReport, err := os.ReadFile(filepath.Join(corpusDir, "expected-report.json"))
	require.NoError(t, err, "missing expected-report.json (run -update)")
	require.JSONEq(t, string(wantReport), normReport, "report mismatch")

	// Idempotency: a second run mutates nothing and proposes no Changed finding.
	before := map[string]string{}
	for _, rel := range gotFiles {
		b, err := os.ReadFile(filepath.Join(tmp, rel))
		require.NoError(t, err)
		before[rel] = string(b)
	}
	report2, _ := runUpgrader(t, tmp)
	for _, rel := range tfFiles(t, tmp) {
		b, err := os.ReadFile(filepath.Join(tmp, rel))
		require.NoError(t, err)
		require.Equal(t, before[rel], string(b), "second run mutated %s", rel)
	}
	require.NotContains(t, severities(t, normalizeReport(report2, tmp)), "Changed", "second run still proposed a rewrite")
}
