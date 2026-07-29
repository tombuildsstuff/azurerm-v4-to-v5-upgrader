data "azurerm_storage_container" "example" {
  name               = "content"
  storage_account_id = azurerm_storage_account.example.id
}
