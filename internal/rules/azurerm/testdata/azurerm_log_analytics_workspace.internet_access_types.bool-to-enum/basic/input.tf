resource "azurerm_log_analytics_workspace" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example-rg"

  internet_ingestion_enabled = true
  internet_query_enabled     = false
}
