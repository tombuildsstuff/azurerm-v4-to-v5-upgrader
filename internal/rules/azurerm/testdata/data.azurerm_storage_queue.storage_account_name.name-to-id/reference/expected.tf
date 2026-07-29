data "azurerm_storage_queue" "example" {
  name               = "content"
  storage_account_id = azurerm_storage_account.example.id
}
