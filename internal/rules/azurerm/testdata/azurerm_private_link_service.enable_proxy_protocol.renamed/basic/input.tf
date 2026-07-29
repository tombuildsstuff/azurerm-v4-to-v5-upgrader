resource "azurerm_private_link_service" "example" {
  name                  = "example"
  resource_group_name   = "example"
  location              = "West Europe"
  enable_proxy_protocol = true

  nat_ip_configuration {
    name      = "primary"
    primary   = true
    subnet_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/virtualNetworks/example/subnets/example"
  }

  load_balancer_frontend_ip_configuration_ids = [
    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/loadBalancers/example/frontendIPConfigurations/example",
  ]
}
