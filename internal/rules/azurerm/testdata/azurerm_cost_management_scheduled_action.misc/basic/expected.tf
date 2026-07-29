resource "azurerm_cost_management_scheduled_action" "example" {
  name    = "example"
  scope   = "example-scope-id"
  view_id = "example-view-id"

  email_subject   = "Monthly Cost Report"
  email_addresses = ["admin@example.com"]

  notification_email {
    subject = "Monthly Cost Report"
    message = "Please find attached the monthly cost report."
  }

  schedule {
    frequency  = "Monthly"
    start_date = "2024-01-01T00:00:00Z"
    end_date   = "2025-01-01T00:00:00Z"
  }
}
