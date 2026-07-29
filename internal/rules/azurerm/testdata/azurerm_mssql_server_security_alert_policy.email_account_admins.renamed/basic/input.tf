resource "azurerm_mssql_server_security_alert_policy" "example" {
  resource_group_name  = "example"
  server_name          = "example"
  state                = "Enabled"
  email_account_admins = true
}
