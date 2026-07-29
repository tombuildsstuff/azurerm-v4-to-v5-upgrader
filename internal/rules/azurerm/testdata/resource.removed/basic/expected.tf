resource "azurerm_app_service" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  app_service_plan_id = "example"
}

resource "azurerm_static_site" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
}

resource "azurerm_postgresql_server" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
}
