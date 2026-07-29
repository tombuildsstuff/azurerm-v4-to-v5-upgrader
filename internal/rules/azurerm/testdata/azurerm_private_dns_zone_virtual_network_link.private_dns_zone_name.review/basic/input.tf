resource "azurerm_private_dns_zone_virtual_network_link" "example" {
  name                  = "example"
  resource_group_name   = "example"
  private_dns_zone_name = "example.com"
  virtual_network_id    = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/virtualNetworks/example"
}
