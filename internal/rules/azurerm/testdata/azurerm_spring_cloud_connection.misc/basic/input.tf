resource "azurerm_spring_cloud_connection" "example" {
  name               = "example"
  spring_cloud_id    = "example-spring-cloud-id"
  target_resource_id = "example-target-resource-id"
  client_type        = "none"

  authentication {
    type = "systemAssignedIdentity"
  }
}
