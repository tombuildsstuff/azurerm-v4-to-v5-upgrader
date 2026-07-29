data "azurerm_key_vault" "example" {
  name                = "example"
  resource_group_name = "example-rg"

  rbac_authorization_enabled = true
}
