resource "azurerm_kusto_eventgrid_data_connection" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  cluster_name        = "example"
  database_name       = "example"
  storage_account_id  = "example-storage-account-id"
  eventhub_id         = "example-eventhub-id"
  consumer_group      = "$Default"

  eventgrid_resource_id        = "example-eventgrid-topic-id"
  managed_identity_resource_id = "example-managed-identity-id"
}
