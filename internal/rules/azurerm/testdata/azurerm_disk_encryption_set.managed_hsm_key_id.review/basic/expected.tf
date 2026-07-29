resource "azurerm_disk_encryption_set" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  managed_hsm_key_id  = "https://example.managedhsm.azure.net/keys/example-key/abcdef0123456789abcdef0123456789"

  identity {
    type = "SystemAssigned"
  }
}
