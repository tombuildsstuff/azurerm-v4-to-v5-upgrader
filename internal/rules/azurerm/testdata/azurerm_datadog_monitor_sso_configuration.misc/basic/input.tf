resource "azurerm_datadog_monitor_sso_configuration" "example" {
  name               = "default"
  datadog_monitor_id = "example-datadog-monitor-id"

  single_sign_on_enabled = "Enable"
}
