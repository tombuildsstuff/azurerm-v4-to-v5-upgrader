data "azurerm_storage_container" "example" {
  name                 = "content"
  storage_account_name = data.azurerm_storage_account.example.name
}
