resource "azurerm_virtual_network_gateway" "example" {
  name                = "example"
  resource_group_name = "example"
  type                = "Vpn"
  bgp_enabled         = true
}
