resource "azurerm_iot_security_solution" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  display_name        = "example"
  iothub_ids          = ["example-iothub-id"]

  recommendations_enabled {
    acr_authentication = true
  }
}
