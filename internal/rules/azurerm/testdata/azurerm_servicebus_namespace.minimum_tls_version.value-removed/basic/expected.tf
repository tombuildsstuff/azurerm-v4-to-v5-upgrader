resource "azurerm_servicebus_namespace" "example" {
  name                = "example-namespace"
  location            = "westeurope"
  resource_group_name = "example-resources"
  sku                 = "Standard"
  minimum_tls_version = "1.0"
}
