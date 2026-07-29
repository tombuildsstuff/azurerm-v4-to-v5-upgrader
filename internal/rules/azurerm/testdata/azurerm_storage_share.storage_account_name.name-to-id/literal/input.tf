resource "azurerm_storage_share" "example" {
  name                 = "example"
  storage_account_name = "mystorageaccount"
  quota                = 50
}
