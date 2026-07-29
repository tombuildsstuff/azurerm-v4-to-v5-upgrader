resource "azurerm_lb_outbound_rule" "example" {
  name                    = "example"
  loadbalancer_id         = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/loadBalancers/example"
  protocol                = "Tcp"
  backend_address_pool_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/loadBalancers/example/backendAddressPools/example"
  enable_tcp_reset        = true

  frontend_ip_configuration {
    name = "example"
  }
}
