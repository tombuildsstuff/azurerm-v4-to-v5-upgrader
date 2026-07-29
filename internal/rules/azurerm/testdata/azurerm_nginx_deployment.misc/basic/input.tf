resource "azurerm_nginx_deployment" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  sku                 = "standard_Monthly"

  diagnose_support_enabled = true
  managed_resource_group   = "example-managed-rg"

  logging_storage_account {
    name           = "examplestorage"
    container_name = "logs"
  }
}
