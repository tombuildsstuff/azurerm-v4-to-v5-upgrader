resource "azurerm_container_app_environment" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"
}
