resource "azurerm_mssql_database" "example" {
  name      = "example"
  server_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Sql/servers/example"
  sku_name  = "S0"

  long_term_retention_policy {
    week_of_year = 3
  }

  threat_detection_policy {
    state                = "Enabled"
    email_account_admins = "Enabled"
    storage_endpoint     = "https://example.blob.core.windows.net/"
  }
}
