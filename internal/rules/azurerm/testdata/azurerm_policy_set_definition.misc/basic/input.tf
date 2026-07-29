resource "azurerm_policy_set_definition" "example" {
  name                = "example"
  policy_type         = "Custom"
  display_name        = "example"
  management_group_id = "example-management-group-id"

  policy_definition_reference {
    policy_definition_id = "example-policy-definition-id"
  }
}
