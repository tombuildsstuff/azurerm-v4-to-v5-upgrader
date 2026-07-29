resource "azurerm_mssql_server_transparent_data_encryption" "example" {
  server_id          = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Sql/servers/example"
  managed_hsm_key_id = "https://example.managedhsm.azure.net/keys/example/abcdef0123456789abcdef0123456789"
}
