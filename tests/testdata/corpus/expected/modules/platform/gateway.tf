# --- Application Gateway ---
# NOTE: authentication_certificate (top-level and inside backend_http_settings)
# is intentionally omitted: it was removed in v5 in favour of
# trusted_root_certificate / trusted_root_certificate_names, and the upgrader
# only flags both occurrences for manual review (flag-only, unrewritten). A v4
# fixture using either would fail GA validate post-upgrade with `Unsupported
# block type: Blocks of type "authentication_certificate" are not expected
# here.` See tests/README.md "Known rule gaps".

resource "azurerm_application_gateway" "example" {
  name                = "example-appgw"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  http2_enabled       = true

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "example"
    subnet_id = azurerm_subnet.appgw.id
  }

  frontend_port {
    name = "port-80"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "example"
    public_ip_address_id = azurerm_public_ip.appgw.id
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
    name                                = "example"
    verify_client_certificate_issuer_dn = true
  }
}
