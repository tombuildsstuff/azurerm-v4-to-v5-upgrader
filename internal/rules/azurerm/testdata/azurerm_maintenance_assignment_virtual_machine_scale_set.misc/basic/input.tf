resource "azurerm_maintenance_assignment_virtual_machine_scale_set" "example" {
  maintenance_configuration_id = "example-maintenance-config-id"
  virtual_machine_scale_set_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/example/providers/microsoft.compute/virtualmachinescalesets/example"
}
