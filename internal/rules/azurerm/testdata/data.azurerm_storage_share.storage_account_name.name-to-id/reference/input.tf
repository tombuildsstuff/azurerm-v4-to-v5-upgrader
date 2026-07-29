data "azurerm_storage_share" "example" {
  name                 = "content"
  storage_account_name = azurerm_storage_account.example.name
}
