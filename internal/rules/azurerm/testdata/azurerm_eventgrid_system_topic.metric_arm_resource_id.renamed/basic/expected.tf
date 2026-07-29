resource "azurerm_eventgrid_system_topic" "example" {
  name                = "example-system-topic"
  location            = "westeurope"
  resource_group_name = "example-resources"
  source_resource_id  = azurerm_storage_account.example.id
  topic_type          = "Microsoft.Storage.StorageAccounts"
  metric_resource_id  = azurerm_storage_account.example.id
}
