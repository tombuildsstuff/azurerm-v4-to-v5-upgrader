data "azurerm_storage_queue" "example" {
  name                 = "content"
  storage_account_name = azurerm_storage_account.example.name
}
