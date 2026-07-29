resource "azurerm_api_management" "example" {
  name = "example"
  security {
    backend_ssl30_enabled = true
  }
}
