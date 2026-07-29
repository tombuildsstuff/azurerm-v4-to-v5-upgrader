resource "azurerm_storage_queue" "example" {
  name                 = "example"
  storage_account_name = "mystorageaccount"
}
