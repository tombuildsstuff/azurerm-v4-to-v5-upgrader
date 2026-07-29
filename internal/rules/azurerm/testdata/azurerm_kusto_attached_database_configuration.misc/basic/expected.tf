resource "azurerm_kusto_attached_database_configuration" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  cluster_name        = "example"
  database_name       = "example"
  cluster_id          = "example-cluster-resource-id"
}
