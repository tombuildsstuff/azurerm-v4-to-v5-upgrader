resource "azurerm_monitor_aad_diagnostic_setting" "example" {
  name                       = "example"
  log_analytics_workspace_id = "example-workspace-id"

  enabled_log {
    category = "SignInLogs"

    retention_policy {
      enabled = true
      days    = 7
    }
  }
}
