resource "azurerm_log_analytics_workspace" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example-rg"

  internet_ingestion_access_type = "Enabled"
  internet_query_access_type     = "Disabled"
}
