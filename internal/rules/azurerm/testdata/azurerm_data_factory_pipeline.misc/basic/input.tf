resource "azurerm_data_factory_pipeline" "example" {
  name            = "example"
  data_factory_id = "example-data-factory-id"

  moniter_metrics_after_duration = "PT1H"
}
