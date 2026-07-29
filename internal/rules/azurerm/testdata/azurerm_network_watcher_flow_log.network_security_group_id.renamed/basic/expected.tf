resource "azurerm_network_watcher_flow_log" "example" {
  network_watcher_name = "example"
  resource_group_name  = "example"
  name                 = "example"

  target_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/networkSecurityGroups/example"
  storage_account_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Storage/storageAccounts/example"
  enabled            = true

  retention_policy {
    enabled = true
    days    = 7
  }
}
