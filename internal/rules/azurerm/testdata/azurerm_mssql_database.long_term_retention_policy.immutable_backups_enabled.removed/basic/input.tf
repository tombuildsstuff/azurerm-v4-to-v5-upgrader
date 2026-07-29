resource "azurerm_mssql_database" "example" {
  name      = "example"
  server_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Sql/servers/example"
  sku_name  = "S0"

  long_term_retention_policy {
    week_of_year              = 5
    immutable_backups_enabled = true
  }
}
