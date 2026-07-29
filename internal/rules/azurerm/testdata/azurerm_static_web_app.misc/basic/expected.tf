resource "azurerm_static_web_app" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = ["example-identity-id"]
  }
}
