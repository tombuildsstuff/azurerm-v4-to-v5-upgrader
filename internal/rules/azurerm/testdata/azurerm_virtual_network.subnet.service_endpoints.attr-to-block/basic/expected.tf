resource "azurerm_virtual_network" "example" {
  name                = "example"
  resource_group_name = "example-rg"
  location            = "westeurope"
  address_space       = ["10.0.0.0/16"]

  subnet {
    name             = "internal"
    address_prefixes = ["10.0.1.0/24"]
    service_endpoint {
      service = "Microsoft.Storage"
    }
  }
}
