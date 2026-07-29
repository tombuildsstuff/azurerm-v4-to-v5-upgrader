data "azurerm_servicebus_subscription" "example" {
  name                = "example-sub"
  topic_name          = "example-topic"
  resource_group_name = "example-rg"
  namespace_name      = "example-ns"
}
