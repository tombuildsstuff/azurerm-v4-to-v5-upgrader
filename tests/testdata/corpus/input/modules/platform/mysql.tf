# --- MySQL Flexible Server ---
# NOTE: customer_managed_key.managed_hsm_key_id is intentionally omitted (not
# just left unset): it's consolidated into key_vault_key_id in v5, and the
# upgrader only flags it for manual review (flag-only, unrewritten). A v4
# fixture using it would fail GA validate post-upgrade with `Unsupported
# argument`. See tests/README.md "Known rule gaps".

resource "azurerm_mysql_flexible_server" "example" {
  name                   = "example-mysql"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_login    = "adminuser"
  administrator_password = "H@Sh1CoR3!"
  sku_name               = "B_Standard_B1s"
  version                = "8.0.21"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  customer_managed_key {
    key_vault_key_id                  = "https://example.vault.azure.net/keys/example/abc123"
    primary_user_assigned_identity_id = azurerm_user_assigned_identity.example.id
  }
}
