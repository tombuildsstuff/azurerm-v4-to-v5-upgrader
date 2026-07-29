resource "azurerm_linux_function_app_slot" "example" {
  name            = "staging"
  function_app_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Web/sites/example-linux-function-app"

  storage_account_name       = "examplestorageacc"
  storage_account_access_key = "storage-access-key-placeholder"

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2017"
  }
}
