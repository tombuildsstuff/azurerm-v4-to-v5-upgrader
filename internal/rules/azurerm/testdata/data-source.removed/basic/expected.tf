data "azurerm_function_app" "example" {
  name                = "example"
  resource_group_name = "example-rg"
}

data "azurerm_postgresql_server" "example" {
  name                = "example"
  resource_group_name = "example-rg"
}
