resource "azurerm_netapp_volume" "example" {
  name                = "example"
  resource_group_name = "example-rg"
  location            = "westeurope"

  mount_ip_addresses = ["10.0.0.4"]
}
