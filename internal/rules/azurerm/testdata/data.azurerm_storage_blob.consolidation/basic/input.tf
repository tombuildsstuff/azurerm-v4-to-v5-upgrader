data "azurerm_storage_blob" "example" {
  name                   = "example.txt"
  storage_account_name   = "myaccount"
  storage_container_name = "content"
}
