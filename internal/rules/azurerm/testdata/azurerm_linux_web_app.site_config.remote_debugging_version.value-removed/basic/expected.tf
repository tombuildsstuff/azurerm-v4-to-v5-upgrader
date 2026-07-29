resource "azurerm_linux_web_app" "example" {
  name                = "example-linux-web-app"
  resource_group_name = "example-resources"
  location            = "West Europe"
  service_plan_id     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Web/serverfarms/example-plan"

  site_config {
    application_stack {
    }

    remote_debugging_enabled = true
    remote_debugging_version = "VS2019"
  }
}
