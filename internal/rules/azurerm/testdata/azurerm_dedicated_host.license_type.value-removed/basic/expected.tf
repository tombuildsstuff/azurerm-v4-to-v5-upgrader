resource "azurerm_dedicated_host_group" "example" {
  name                        = "example-host-group"
  resource_group_name         = "example-resources"
  location                    = "West Europe"
  platform_fault_domain_count = 2
}

resource "azurerm_dedicated_host" "example" {
  name                    = "example-host"
  location                = "West Europe"
  dedicated_host_group_id = azurerm_dedicated_host_group.example.id
  sku_name                = "DSv3-Type3"
  platform_fault_domain   = 1
  license_type            = "None"
}
