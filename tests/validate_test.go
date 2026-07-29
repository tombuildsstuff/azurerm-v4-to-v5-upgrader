package tests

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type tfDiag struct {
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
	Detail   string `json:"detail"`
}

// terraformValidateErrors runs `terraform validate -json` in dir and returns
// only error-severity diagnostics. env carries TF_CLI_CONFIG_FILE (dev_overrides
// disabled) and TF_PLUGIN_CACHE_DIR so the registry provider is used and cached.
func terraformValidateErrors(t *testing.T, dir string, env []string) []tfDiag {
	t.Helper()
	cmd := exec.Command("terraform", "-chdir="+dir, "validate", "-json")
	cmd.Env = env
	out, _ := cmd.CombinedOutput()
	var parsed struct {
		Diagnostics []tfDiag `json:"diagnostics"`
	}
	require.NoError(t, json.Unmarshal(out, &parsed), "validate -json not JSON: %s", out)
	var errs []tfDiag
	for _, d := range parsed.Diagnostics {
		if d.Severity == "error" {
			errs = append(errs, d)
		}
	}
	return errs
}

func TestValidate(t *testing.T) {
	if os.Getenv("AZURERM_ACC") == "" {
		t.Skip("acceptance: set AZURERM_ACC=1 (needs terraform + network for the registry)")
	}
	if _, err := exec.LookPath("terraform"); err != nil {
		t.Skip("terraform not on PATH")
	}

	tmp := t.TempDir()
	copyTree(t, filepath.Join(corpusDir, "input"), tmp)
	runUpgrader(t, tmp) // upgrade in place, then validate the result

	// Validate against the RELEASED registry provider, not any dev_overrides:
	// write a CLI config that forces direct (registry) installation.
	rc := filepath.Join(tmp, "cli.tfrc")
	require.NoError(t, os.WriteFile(rc, []byte("provider_installation {\n  direct {}\n}\n"), 0o644))

	// Persistent plugin cache so the ~340MB provider is downloaded once.
	cache := os.Getenv("TF_PLUGIN_CACHE_DIR")
	if cache == "" {
		cache = filepath.Join(os.TempDir(), "upgrader-tf-plugin-cache")
	}
	require.NoError(t, os.MkdirAll(cache, 0o755))

	env := append(os.Environ(),
		"TF_CLI_CONFIG_FILE="+rc,
		"TF_PLUGIN_CACHE_DIR="+cache,
		"TF_IN_AUTOMATION=1",
	)

	// init installs the local child-module manifest AND the pinned registry
	// provider; validate alone cannot resolve local modules. On a cold plugin
	// cache the registry download can intermittently hit a transient network
	// error, so retry init a bounded number of times; a deterministic failure
	// (bad provider constraint, malformed upgraded HCL) will still fail every
	// attempt and end in t.Fatalf.
	const maxInitAttempts = 3
	var out []byte
	var err error
	for attempt := 1; attempt <= maxInitAttempts; attempt++ {
		initCmd := exec.Command("terraform", "-chdir="+tmp, "init", "-input=false", "-no-color")
		initCmd.Env = env
		out, err = initCmd.CombinedOutput()
		if err == nil {
			break
		}
		if attempt == maxInitAttempts {
			t.Fatalf("terraform init failed after %d attempts: %v\n%s", maxInitAttempts, err, out)
		}
		time.Sleep(time.Duration(attempt) * 2 * time.Second)
	}

	errs := terraformValidateErrors(t, tmp, env)
	require.Empty(t, errs, "upgraded corpus has v5 validate errors: %+v", errs)
}
