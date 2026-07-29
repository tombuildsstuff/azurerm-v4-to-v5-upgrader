data "azurerm_servicebus_namespace_disaster_recovery_config" "example" {
  name                = "example-alias"
  resource_group_name = "example-rg"
  namespace_name      = "example-ns"
}
