resource "azurerm_log_analytics_workspace" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"
  sku                 = "PerGB2018"

  local_authentication_disabled = true
}
