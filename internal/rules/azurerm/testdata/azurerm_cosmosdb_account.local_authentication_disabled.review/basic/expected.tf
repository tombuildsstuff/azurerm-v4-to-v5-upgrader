resource "azurerm_cosmosdb_account" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  local_authentication_disabled = true
  managed_hsm_key_id            = "https://example.managedhsm.azure.net/keys/example/abcdef0123456789abcdef0123456789"
  minimal_tls_version           = "Tls11"

  consistency_policy {
    consistency_level = "Session"
  }

  geo_location {
    location          = "westeurope"
    failover_priority = 0
  }
}
