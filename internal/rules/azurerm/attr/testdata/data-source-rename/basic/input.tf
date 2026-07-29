data "azurerm_key_vault" "example" {
  name                = "example"
  resource_group_name = "example-rg"

  enable_rbac_authorization = true
}
