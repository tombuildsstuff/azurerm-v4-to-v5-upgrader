resource "azurerm_api_management" "example" {
  name = "example"
  security {
    enable_backend_ssl30 = true
  }
}
