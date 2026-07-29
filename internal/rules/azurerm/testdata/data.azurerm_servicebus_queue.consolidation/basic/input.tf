data "azurerm_servicebus_queue" "example" {
  name                = "example-queue"
  resource_group_name = "example-rg"
  namespace_name      = "example-ns"
}
