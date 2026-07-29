resource "azurerm_storage_share_directory" "example" {
  name             = "example"
  storage_share_id = azurerm_storage_share.example.id
}
