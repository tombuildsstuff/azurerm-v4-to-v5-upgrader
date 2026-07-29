resource "azurerm_key_vault" "example" {
  name                      = "example"
  location                  = "westeurope"
  resource_group_name       = "example"
  tenant_id                 = "00000000-0000-0000-0000-000000000000"
  sku_name                  = "standard"
  enable_rbac_authorization = true
}
