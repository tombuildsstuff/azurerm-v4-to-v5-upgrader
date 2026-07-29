resource "azurerm_federated_identity_credential" "example" {
  name                      = "example"
  resource_group_name       = "example"
  audience                  = ["api://AzureADTokenExchange"]
  issuer                    = "https://example.com/issuer"
  user_assigned_identity_id = azurerm_user_assigned_identity.example.id
  subject                   = "system:serviceaccount:example:example"
}
