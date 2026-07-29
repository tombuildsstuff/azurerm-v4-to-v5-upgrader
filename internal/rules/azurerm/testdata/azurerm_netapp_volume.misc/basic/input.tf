resource "azurerm_netapp_volume" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"
  account_name        = "example"
  pool_name           = "example"
  volume_path         = "example"
  service_level       = "Standard"
  subnet_id           = "example-subnet-id"
  storage_quota_in_gb = 100

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }
}
