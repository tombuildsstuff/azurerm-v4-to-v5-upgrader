# --- Cost Management scheduled action ---
# (email_subject is flagged NeedsReview in v5 because it now validates a
# maximum length of 50 characters; the value below already fits, so this is
# purely informational and does not affect GA validate. NOTE: the
# azurerm_cost_management_scheduled_action.misc per-rule fixture's `scope`
# argument and nested `notification_email`/`schedule` blocks reflect an older
# schema shape; GA v5.0.0's actual schema uses top-level
# frequency/email_address_sender/display_name/start_date/end_date instead, as
# used below - this mismatch predates v5 (not itself a v4->v5 change) so it's
# just a fixture-completeness fix, not a rule gap.)

resource "azurerm_cost_management_scheduled_action" "example" {
  name                 = "example-cost-action"
  view_id              = "${azurerm_resource_group.example.id}/providers/Microsoft.CostManagement/views/ms:DailyCosts"
  display_name         = "Monthly Cost Report"
  frequency            = "Monthly"
  start_date           = "2024-01-01T00:00:00Z"
  end_date             = "2025-01-01T00:00:00Z"
  email_subject        = "Monthly Cost Report"
  email_addresses      = ["admin@example.com"]
  email_address_sender = "admin@example.com"
}
