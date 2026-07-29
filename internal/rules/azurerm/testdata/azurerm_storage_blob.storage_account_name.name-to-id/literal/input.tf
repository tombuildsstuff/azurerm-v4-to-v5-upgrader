resource "azurerm_storage_blob" "example" {
  name                   = "example.vhd"
  storage_account_name   = "mystorageaccount"
  storage_container_name = "example"
  type                   = "Page"
  source                 = "some-local-file.vhd"
}
