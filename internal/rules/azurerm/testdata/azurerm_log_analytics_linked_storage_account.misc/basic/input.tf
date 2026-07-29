resource "azurerm_log_analytics_linked_storage_account" "example" {
  data_source_type      = "CustomLogs"
  workspace_resource_id = "example-workspace-id"
  storage_account_ids   = ["example-storage-account-id"]
}
