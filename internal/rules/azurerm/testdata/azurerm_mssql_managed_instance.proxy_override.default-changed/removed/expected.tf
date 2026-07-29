resource "azurerm_mssql_managed_instance" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  license_type        = "LicenseIncluded"
  sku_name            = "GP_Gen5"
  storage_size_in_gb  = 32
  subnet_id           = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/virtualNetworks/example/subnets/example"
  vcores              = 4

  administrator_login          = "exampleadmin"
  administrator_login_password = "ExamplePassword123!"

  minimum_tls_version = "1.0"
  proxy_override      = "Default"
}
