resource "azurerm_recovery_services_vault" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  sku                 = "Standard"
}
