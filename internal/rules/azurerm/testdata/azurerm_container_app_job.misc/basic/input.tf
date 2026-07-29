resource "azurerm_container_app_job" "example" {
  name                         = "example"
  container_app_environment_id = "example-env-id"
  resource_group_name          = "example"
  replica_timeout_in_seconds   = 10
  trigger_type                 = "Manual"

  template {
    container {
      name   = "example"
      image  = "nginx:latest"
      cpu    = 0.25
      memory = "0.5Gi"

      liveness_probe {
        transport                        = "HTTP"
        port                             = 80
        termination_grace_period_seconds = 10
      }

      startup_probe {
        transport                        = "HTTP"
        port                             = 80
        termination_grace_period_seconds = 10
      }
    }
  }
}
