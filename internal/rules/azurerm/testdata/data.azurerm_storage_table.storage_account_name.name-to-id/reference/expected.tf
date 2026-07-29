data "azurerm_storage_table" "example" {
  name               = "content"
  storage_account_id = azurerm_storage_account.example.id
}
