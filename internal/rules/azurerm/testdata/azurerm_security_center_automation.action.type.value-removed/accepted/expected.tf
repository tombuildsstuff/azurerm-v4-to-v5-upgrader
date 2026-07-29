resource "azurerm_security_center_automation" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  scopes              = ["/subscriptions/00000000-0000-0000-0000-000000000000"]

  action {
    type        = "LogicApp"
    resource_id = azurerm_logic_app_workflow.example.id
    trigger_url = "https://example.com/trigger"
  }

  source {
    event_source = "Alerts"
  }
}
