resource "azurerm_batch_pool" "example" {
  name                = "example"
  resource_group_name = "example"
  account_name        = "example"
  display_name        = "example"
  vm_size             = "Standard_A1"

  certificate {
    id             = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Batch/batchAccounts/example/certificates/sha1-abcdef"
    store_location = "CurrentUser"
  }
}
