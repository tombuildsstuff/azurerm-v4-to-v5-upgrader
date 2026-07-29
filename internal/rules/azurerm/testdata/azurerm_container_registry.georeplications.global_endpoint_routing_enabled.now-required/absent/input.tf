resource "azurerm_container_registry" "example" {
  name                = "example"
  resource_group_name = "example-rg"
  location            = "westeurope"
  sku                 = "Premium"

  georeplications {
    location = "northeurope"
  }
}
