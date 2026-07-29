# --- Automation account ---
# (encryption.key_source is removed in v5; the upgrader strips it. Since a
# bare `key_source = "Microsoft.Automation"` (Microsoft-managed keys, the
# default) would leave an empty `encryption {}` block that GA v5.0.0 rejects
# (key_vault_key_id becomes required whenever the block is present), this
# uses a realistic CMK config instead - key_vault_key_id and
# user_assigned_identity_id were already valid v4 fields alongside
# key_source, so this is not new v5-only syntax.)

resource "azurerm_automation_account" "example" {
  name                = "example-automation"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  encryption {
    key_source                = "Microsoft.Keyvault"
    key_vault_key_id          = "https://example.vault.azure.net/keys/example/abc123"
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
  }
}
