resource "azurerm_express_route_connection" "example" {
  name                             = "example"
  express_route_circuit_peering_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/expressRouteCircuits/example/peerings/AzurePrivatePeering"
  express_route_gateway_id         = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/expressRouteGateways/example"
  internet_security_enabled        = true
}
