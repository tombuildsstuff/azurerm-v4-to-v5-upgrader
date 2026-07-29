resource "azurerm_container_app" "example" {
  name                         = "example"
  container_app_environment_id = "example-env-id"
  resource_group_name          = "example"
  revision_mode                = "Single"

  template {
    container {
      name   = "example"
      image  = "nginx:latest"
      cpu    = 0.25
      memory = "0.5Gi"

      liveness_probe {
        transport = "HTTP"
        port      = 80
      }

      startup_probe {
        transport = "HTTP"
        port      = 80
      }
    }
  }
}
