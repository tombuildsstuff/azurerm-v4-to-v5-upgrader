resource "azurerm_application_insights" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"
  application_type    = "web"

  daily_data_cap_notifications_disabled = true
  disable_ip_masking                    = true
  local_authentication_disabled         = true
}
