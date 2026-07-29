resource "azurerm_application_gateway" "example" {
  name                = "example"
  resource_group_name = "example-rg"
  location            = "westeurope"

  authentication_certificate {
    name = "cert"
    data = "PEM"
  }

  backend_http_settings {
    name = "settings"
    authentication_certificate {
      name = "cert"
    }
  }
}
