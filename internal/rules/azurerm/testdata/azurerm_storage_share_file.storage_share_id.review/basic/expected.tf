resource "azurerm_storage_share_file" "example" {
  name             = "example.txt"
  storage_share_id = azurerm_storage_share.example.id
  source           = "some-local-file.txt"
}
