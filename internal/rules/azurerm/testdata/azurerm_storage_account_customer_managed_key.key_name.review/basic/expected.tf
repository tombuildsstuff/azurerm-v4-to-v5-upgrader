resource "azurerm_storage_account_customer_managed_key" "example" {
  storage_account_id = azurerm_storage_account.example.id
  key_vault_uri      = azurerm_key_vault.example.vault_uri
  key_name           = "example"
  key_version        = "1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d"
}
