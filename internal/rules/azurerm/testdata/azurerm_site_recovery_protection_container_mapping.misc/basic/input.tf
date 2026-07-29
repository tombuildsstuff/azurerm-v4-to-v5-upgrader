resource "azurerm_site_recovery_protection_container_mapping" "example" {
  name                                      = "example"
  resource_group_name                       = "example"
  recovery_vault_name                       = "example"
  recovery_fabric_name                      = "example"
  recovery_source_protection_container_name = "example"
  recovery_target_protection_container_id   = "example-target-container-id"
  recovery_replication_policy_id            = "example-policy-id"

  automatic_update {
    enabled = true
  }
}
