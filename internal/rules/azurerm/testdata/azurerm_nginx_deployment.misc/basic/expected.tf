resource "azurerm_nginx_deployment" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  sku                 = "standard_Monthly"


  logging_storage_account {
    name           = "examplestorage"
    container_name = "logs"
  }
}
