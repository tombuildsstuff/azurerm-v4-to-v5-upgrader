resource "azurerm_storage_account_customer_managed_key" "example" {
  storage_account_id = azurerm_storage_account.example.id
  managed_hsm_key_id = azurerm_key_vault_managed_hardware_security_module_key.example.id
}
