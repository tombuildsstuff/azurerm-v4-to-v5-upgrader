resource "azurerm_dashboard_grafana" "example" {
  name                  = "example"
  resource_group_name   = "example"
  location              = "westeurope"
  sku                   = "Essential"
  grafana_major_version = "11"
}
