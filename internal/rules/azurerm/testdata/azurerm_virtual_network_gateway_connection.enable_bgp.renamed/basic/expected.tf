resource "azurerm_virtual_network_gateway_connection" "example" {
  name                       = "example"
  location                   = "West Europe"
  resource_group_name        = "example"
  type                       = "IPsec"
  virtual_network_gateway_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/virtualNetworkGateways/example"
  local_network_gateway_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/localNetworkGateways/example"
  shared_key                 = "changeme"
  bgp_enabled                = true
}
