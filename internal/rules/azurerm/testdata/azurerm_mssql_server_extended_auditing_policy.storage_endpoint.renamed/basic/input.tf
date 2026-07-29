resource "azurerm_mssql_server_extended_auditing_policy" "example" {
  server_id                  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Sql/servers/example"
  storage_endpoint           = "https://example.blob.core.windows.net/"
  storage_account_access_key = "example-access-key"
  retention_in_days          = 6
}
