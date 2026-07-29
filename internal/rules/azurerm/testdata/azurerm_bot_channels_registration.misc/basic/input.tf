resource "azurerm_bot_channels_registration" "example" {
  name                = "example"
  location            = "global"
  resource_group_name = "example"
  sku                 = "F0"
  microsoft_app_id    = "00000000-0000-0000-0000-000000000000"
}
