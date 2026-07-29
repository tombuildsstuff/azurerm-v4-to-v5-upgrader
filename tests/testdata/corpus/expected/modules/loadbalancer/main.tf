resource "azurerm_lb_nat_rule" "example" {
  resource_group_name            = "example"
  loadbalancer_id                = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/loadBalancers/example"
  name                           = "example"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "example"
  floating_ip_enabled            = true
  tcp_reset_enabled              = true
}
