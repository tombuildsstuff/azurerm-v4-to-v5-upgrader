resource "azurerm_bot_channel_ms_teams" "example" {
  bot_name            = "example"
  location            = "global"
  resource_group_name = "example"

  enable_calling = true
}
