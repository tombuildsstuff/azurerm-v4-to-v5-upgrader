resource "azurerm_linux_web_app_slot" "example" {
  name           = "staging"
  app_service_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Web/sites/example-linux-web-app"

  site_config {
    application_stack {
      ruby_version = "2.6"
    }

    remote_debugging_enabled = true
    remote_debugging_version = "VS2017"
  }
}
