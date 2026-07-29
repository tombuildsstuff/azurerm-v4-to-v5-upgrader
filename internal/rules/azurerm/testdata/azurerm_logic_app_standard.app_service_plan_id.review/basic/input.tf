resource "azurerm_logic_app_standard" "example" {
  name                       = "example-logic-app"
  resource_group_name        = "example-resources"
  location                   = "West Europe"
  app_service_plan_id        = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/example-resources/providers/microsoft.web/serverfarms/example-plan"
  storage_account_name       = "examplestorageacc"
  storage_account_access_key = "storage-access-key-placeholder"

  site_config {
    min_tls_version     = "1.1"
    scm_min_tls_version = "1.0"

    public_network_access_enabled = true
  }
}
