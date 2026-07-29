resource "azurerm_eventhub" "example" {
  name                = "example-eventhub"
  namespace_name      = "example-namespace"
  resource_group_name = "example-resources"
  partition_count     = 2
  message_retention   = 1
}
