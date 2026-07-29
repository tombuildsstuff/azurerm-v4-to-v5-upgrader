resource "azurerm_mysql_flexible_server" "example" {
  name                   = "example"
  resource_group_name    = "example"
  location               = "westeurope"
  administrator_login    = "adminuser"
  administrator_password = "H@Sh1CoR3!"
  sku_name               = "B_Standard_B1s"

  customer_managed_key {
    key_vault_key_id   = "https://example.vault.azure.net/keys/example/abc123"
    managed_hsm_key_id = "https://example.managedhsm.azure.net/keys/example/abc123"
  }
}
