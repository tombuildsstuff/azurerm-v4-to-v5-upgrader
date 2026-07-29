# --- User-assigned identity ---
# (used below by azurerm_mysql_flexible_server's customer_managed_key.)
#
# NOTE: azurerm_federated_identity_credential is intentionally NOT included
# in this module. Its only v5 rule (parent_id -> user_assigned_identity_id)
# renames the reference attribute but does not also strip the now-obsolete
# resource_group_name, which a genuine v4 config always sets (it's required
# in v4). GA v5.0.0 rejects the upgraded output with `Unsupported argument:
# An argument named "resource_group_name" is not expected here.` This is a
# genuine rule bug (the rule needs to also remove resource_group_name), not a
# fixture-completeness issue. See tests/README.md "Known rule gaps".

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-uai"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
