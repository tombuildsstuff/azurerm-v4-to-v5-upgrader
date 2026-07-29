resource "azurerm_mssql_server" "example" {
  name                         = "example"
  resource_group_name          = "example"
  location                     = "westeurope"
  version                      = "12.0"
  administrator_login          = "exampleadmin"
  administrator_login_password = "ExamplePassword123!"
  minimum_tls_version          = "Disabled"
}
