resource "azurerm_mssql_managed_database" "example" {
  name                = "example"
  managed_instance_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Sql/managedInstances/example"

  long_term_retention_policy {
    week_of_year              = 5
    immutable_backups_enabled = true
  }

  short_term_retention_days = 7
}
