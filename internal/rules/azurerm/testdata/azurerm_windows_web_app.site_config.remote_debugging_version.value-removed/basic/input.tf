resource "azurerm_windows_web_app" "example" {
  name                = "example-windows-web-app"
  resource_group_name = "example-resources"
  location            = "West Europe"
  service_plan_id     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Web/serverfarms/example-plan"

  virtual_network_subnet_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Network/virtualNetworks/example-vnet/subnets/example-subnet"

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2019"
  }
}
