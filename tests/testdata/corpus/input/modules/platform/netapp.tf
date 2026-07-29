# --- NetApp account / pool / volume ---
# NOTE: mount_ip_addresses (per the
# azurerm_netapp_volume.mount_ip_addresses.review per-rule fixture) is not
# exercised here: it is a Computed-only attribute (the mount IPs Azure
# assigns), not one a genuine v4 config sets, so there is no realistic v4
# usage to include.

resource "azurerm_netapp_account" "example" {
  name                = "example-netapp"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_pool" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "example"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.netapp.id
  storage_quota_in_gb = 100

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }
}
