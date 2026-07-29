package engine

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// tfDiag is one entry of `terraform validate -json`'s "diagnostics" array.
type tfDiag struct {
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
	Detail   string `json:"detail"`
	Snippet  struct {
		Code string `json:"code"`
	} `json:"snippet"`
}

// runTerraformValidate writes cfg to a temp dir and returns the diagnostics
// from `terraform validate -json`. It fails the test if terraform can't be run
// at all (the caller gates on availability first).
func runTerraformValidate(t *testing.T, cfg string) []tfDiag {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(cfg), 0o644))
	out, _ := exec.Command("terraform", "-chdir="+dir, "validate", "-json").CombinedOutput()
	var parsed struct {
		Diagnostics []tfDiag `json:"diagnostics"`
	}
	// terraform validate exits non-zero when there are errors; the JSON is
	// still on stdout, so parse regardless of exit code.
	require.NoError(t, json.Unmarshal(out, &parsed), "terraform validate -json output was not JSON: %s", out)
	return parsed.Diagnostics
}

// mentions reports whether any error- or warning-severity diagnostic
// references token, in its summary, detail, or source snippet. Both
// severities are considered: against the real v5.4.x provider, several
// v4-only fields are currently surfaced as deprecation *warnings* rather
// than hard errors (the provider keeps them working for a transition
// period ahead of a later removal), so restricting to "error" would miss
// those. The snippet is checked too because deprecation-warning messages
// often name only the *replacement* attribute in summary/detail (e.g.
// "This property has been renamed to `bgp_enabled`") without repeating the
// old token - the old token is only visible in the offending source line
// itself (the snippet). The only other diagnostic this test ever sees is
// the fixed "Provider development overrides are in effect" notice (from
// ~/.terraformrc dev_overrides), which has no snippet and never mentions a
// resource attribute name, so widening the match introduces no risk of a
// false-positive on that notice.
func mentions(diags []tfDiag, token string) bool {
	for _, d := range diags {
		if d.Severity != "error" && d.Severity != "warning" {
			continue
		}
		if strings.Contains(d.Summary, token) || strings.Contains(d.Detail, token) || strings.Contains(d.Snippet.Code, token) {
			return true
		}
	}
	return false
}

const accProvider = `provider "azurerm" {
  features {}
  resource_provider_registrations = "none"
}
`

// TestAcceptanceV5Validate proves, against the real azurerm v5 provider schema
// (via terraform validate), that the upgrader removes v4-only constructs the
// provider rejects. Opt-in: set AZURERM_ACC=1 and have terraform + a v5 azurerm
// provider available (the repo's dev_overrides, or a real install).
func TestAcceptanceV5Validate(t *testing.T) {
	if os.Getenv("AZURERM_ACC") == "" {
		t.Skip("acceptance test: set AZURERM_ACC=1 (needs terraform + an azurerm v5 provider) to run")
	}
	if _, err := exec.LookPath("terraform"); err != nil {
		t.Skip("acceptance test: terraform not on PATH")
	}
	// Probe: terraform validate must run at the CLI level on a trivial config.
	probe := runTerraformValidate(t, accProvider)
	_ = probe // if this panicked/failed, the require inside would have stopped us

	cases := []struct {
		name  string
		token string // the v4 field name that v5 rejects
		cfg   string // a minimal v4 config exhibiting the change (provider block prepended)
	}{
		{
			name:  "virtual_network_gateway_enable_bgp",
			token: "enable_bgp",
			cfg: `resource "azurerm_virtual_network_gateway" "e" {
  name                = "example-gw"
  resource_group_name = "example-rg"
  location            = "westeurope"
  type                = "Vpn"
  sku                 = "VpnGw1"
  enable_bgp          = true
  ip_configuration {
    public_ip_address_id          = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/publicIPAddresses/pip"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/GatewaySubnet"
  }
}
`,
		},
		{
			name:  "key_vault_enable_rbac_authorization",
			token: "enable_rbac_authorization",
			cfg: `resource "azurerm_key_vault" "e" {
  name                      = "example-kv"
  resource_group_name       = "example-rg"
  location                  = "westeurope"
  tenant_id                 = "00000000-0000-0000-0000-000000000000"
  sku_name                  = "standard"
  enable_rbac_authorization = true
}
`,
		},
		{
			name:  "lb_nat_rule_enable_floating_ip",
			token: "enable_floating_ip",
			cfg: `resource "azurerm_lb_nat_rule" "e" {
  name                           = "example-nat"
  resource_group_name            = "example-rg"
  loadbalancer_id                = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/loadBalancers/lb"
  frontend_ip_configuration_name = "feip"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  enable_floating_ip             = true
}
`,
		},
		{
			// v5 makes template.container.liveness_probe.termination_grace_period_seconds
			// computed-only (it's decided automatically post-apply); setting it
			// in config is a hard "Value for unconfigurable attribute" error,
			// not just a deprecation warning.
			name:  "container_app_liveness_probe_termination_grace_period_seconds",
			token: "termination_grace_period_seconds",
			cfg: `resource "azurerm_container_app" "e" {
  name                         = "example-app"
  container_app_environment_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.App/managedEnvironments/example-env"
  resource_group_name          = "example-rg"
  revision_mode                = "Single"
  template {
    container {
      name   = "example"
      image  = "nginx:latest"
      cpu    = 0.25
      memory = "0.5Gi"
      liveness_probe {
        transport                        = "HTTP"
        port                             = 80
        termination_grace_period_seconds = 30
      }
    }
  }
}
`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			full := accProvider + tc.cfg

			// (1) sanity: the v4 config is flagged by the v5 provider on the token.
			v4diags := runTerraformValidate(t, full)
			require.True(t, mentions(v4diags, tc.token),
				"expected v5 provider to flag v4 token %q in the input; diagnostics: %+v", tc.token, v4diags)

			// (2) run the upgrader in place.
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(full), 0o644))
			_, err := Run(Options{Dir: dir})
			require.NoError(t, err)
			out, err := os.ReadFile(filepath.Join(dir, "main.tf"))
			require.NoError(t, err)

			// (3) the transformed output no longer trips the token.
			v5diags := runTerraformValidate(t, string(out))
			require.False(t, mentions(v5diags, tc.token),
				"upgrader output still trips v4 token %q; diagnostics: %+v\noutput:\n%s", tc.token, v5diags, out)
		})
	}
}
