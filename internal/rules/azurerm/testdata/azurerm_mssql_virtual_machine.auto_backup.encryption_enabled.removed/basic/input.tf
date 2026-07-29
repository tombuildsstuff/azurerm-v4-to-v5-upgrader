resource "azurerm_mssql_virtual_machine" "example" {
  virtual_machine_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Compute/virtualMachines/example"
  sql_license_type   = "PAYG"

  auto_backup {
    retention_period_in_days   = 30
    storage_blob_endpoint      = "https://example.blob.core.windows.net/"
    storage_account_access_key = "example-access-key"
    encryption_enabled         = true
    encryption_password        = "ExamplePassword123!"
  }
}
