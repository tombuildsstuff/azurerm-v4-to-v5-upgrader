# --- Policy set definition ---
# (management_group_id is removed entirely in v5; the upgrader strips it, so
# this becomes a subscription-scoped custom policy set definition.)

resource "azurerm_policy_set_definition" "example" {
  name                = "example-policyset"
  policy_type         = "Custom"
  display_name        = "example"
  management_group_id = "example-management-group-id"

  policy_definition_reference {
    policy_definition_id = "/providers/Microsoft.Authorization/policyDefinitions/06a78e20-9358-41c9-923c-fb736d382a4d"
  }
}
