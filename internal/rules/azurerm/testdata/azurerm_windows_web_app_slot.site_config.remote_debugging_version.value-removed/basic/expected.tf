resource "azurerm_windows_web_app_slot" "example" {
  name           = "staging"
  app_service_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Web/sites/example-windows-web-app"

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2017"
  }
}
