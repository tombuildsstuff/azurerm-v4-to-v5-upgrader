resource "azurerm_subnet" "example" {
  name                 = "example"
  resource_group_name  = "example-rg"
  virtual_network_name = "example-vnet"
  address_prefixes     = ["10.0.1.0/24"]

  service_endpoint {
    service = "Microsoft.Storage"
  }
  service_endpoint {
    service = "Microsoft.Sql"
  }
}
