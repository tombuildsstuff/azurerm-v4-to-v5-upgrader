resource "azurerm_monitor_diagnostic_setting" "example" {
  name                       = "example"
  target_resource_id         = "example-target-resource-id"
  log_analytics_workspace_id = "example-workspace-id"

  enabled_log {
    category = "AuditEvent"

    retention_policy {
      enabled = true
      days    = 7
    }
  }

  metric {
    category = "AllMetrics"
    enabled  = true

    retention_policy {
      enabled = true
      days    = 7
    }
  }
}
