resource "azurerm_container_registry" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  sku                 = "Premium"

}
