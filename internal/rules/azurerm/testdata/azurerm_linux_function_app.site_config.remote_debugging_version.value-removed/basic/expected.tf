resource "azurerm_linux_function_app" "example" {
  name                = "example-linux-function-app"
  resource_group_name = "example-resources"
  location            = "West Europe"
  service_plan_id     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Web/serverfarms/example-plan"

  storage_account_name       = "examplestorageacc"
  storage_account_access_key = "storage-access-key-placeholder"

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2019"
  }
}
