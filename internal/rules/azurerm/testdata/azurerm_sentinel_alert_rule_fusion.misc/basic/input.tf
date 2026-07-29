resource "azurerm_sentinel_alert_rule_fusion" "example" {
  name                       = "example"
  log_analytics_workspace_id = "example-workspace-id"
  alert_rule_template_guid   = "00000000-0000-0000-0000-000000000000"
}
