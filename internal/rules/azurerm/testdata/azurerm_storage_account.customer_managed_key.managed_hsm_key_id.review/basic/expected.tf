resource "azurerm_storage_account" "example" {
  name                            = "example"
  resource_group_name             = "example"
  location                        = "westeurope"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
  min_tls_version                 = "TLS1_2"

  identity {
    type = "SystemAssigned"
  }

  customer_managed_key {
    managed_hsm_key_id        = azurerm_key_vault_managed_hardware_security_module_key.example.id
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
  }
}
