resource "azurerm_data_factory_linked_service_azure_databricks" "example" {
  name             = "example"
  data_factory_id  = "example-data-factory-id"
  adb_domain       = "https://example.azuredatabricks.net"
  msi_workspace_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Databricks/workspaces/example"
}
