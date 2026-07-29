resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = "example"
  location                 = "westeurope"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  min_tls_version          = "TLS1_0"
}
