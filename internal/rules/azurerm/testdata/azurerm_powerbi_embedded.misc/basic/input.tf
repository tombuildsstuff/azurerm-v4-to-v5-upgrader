resource "azurerm_powerbi_embedded" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  sku_name            = "A1"
  administrators      = ["admin@example.com"]
}
