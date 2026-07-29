resource "azurerm_mssql_database_extended_auditing_policy" "example" {
  database_id                             = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Sql/servers/example/databases/example"
  blob_storage_endpoint                   = "https://example.blob.core.windows.net/"
  storage_account_access_key              = "example-access-key"
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 6
}
