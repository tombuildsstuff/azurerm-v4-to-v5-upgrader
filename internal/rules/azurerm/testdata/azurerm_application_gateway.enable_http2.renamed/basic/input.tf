resource "azurerm_application_gateway" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"
  enable_http2        = true

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "example"
    subnet_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/virtualNetworks/example/subnets/example"
  }

  frontend_port {
    name = "port-80"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "example"
    public_ip_address_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/publicIPAddresses/example"
  }

  backend_address_pool {
    name = "example"
  }

  backend_http_settings {
    name                  = "example"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 60
  }

  http_listener {
    name                           = "example"
    frontend_ip_configuration_name = "example"
    frontend_port_name             = "port-80"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "example"
    priority                   = 100
    rule_type                  = "Basic"
    http_listener_name         = "example"
    backend_address_pool_name  = "example"
    backend_http_settings_name = "example"
  }

  ssl_profile {
    name                         = "example"
    verify_client_cert_issuer_dn = true
  }
}
