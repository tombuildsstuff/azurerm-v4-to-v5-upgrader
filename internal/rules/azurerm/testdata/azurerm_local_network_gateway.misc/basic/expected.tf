resource "azurerm_local_network_gateway" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"
  gateway_address     = "10.0.0.1"
  address_space       = ["10.1.0.0/16"]
}
