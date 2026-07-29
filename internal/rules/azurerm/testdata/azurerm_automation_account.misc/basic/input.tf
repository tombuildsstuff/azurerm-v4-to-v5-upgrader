resource "azurerm_automation_account" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"
  sku_name            = "Basic"

  encryption {
    key_source = "Microsoft.Automation"
  }
}
