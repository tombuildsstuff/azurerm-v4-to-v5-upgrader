resource "azurerm_data_factory_integration_runtime_self_hosted" "example" {
  name            = "example"
  data_factory_id = "example-data-factory-id"

  rbac_authorization {
    resource_id = "example-integration-runtime-resource-id"
  }
}
