resource "azurerm_eventhub_namespace" "example" {
  name                = "example-namespace"
  location            = "westeurope"
  resource_group_name = "example-resources"
  sku                 = "Standard"
  capacity            = 1
  minimum_tls_version = "1.1"
}
