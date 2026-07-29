# --- Cosmos DB account ---
# NOTE: local_authentication_disabled and managed_hsm_key_id are intentionally
# omitted here (not just left unset): both v5 changes are flag-only -
# local_authentication_disabled's replacement (local_authentication_enabled)
# inverts the boolean meaning, and managed_hsm_key_id is consolidated into
# key_vault_key_id - so the upgrader only flags them for manual review
# without rewriting. A v4 fixture using either would fail GA validate
# post-upgrade with `Unsupported argument`. minimal_tls_version uses "Tls12"
# (a still-accepted value) rather than the removed "Tls11"/"Tls" values the
# azurerm_cosmosdb_account.local_authentication_disabled.review per-rule
# fixture exercises, per the value-removed-category convention (see
# azurerm_mssql_server.minimum_tls_version in modules/mssql). See
# tests/README.md "Known rule gaps".

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-cosmos"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  minimal_tls_version = "Tls12"

  consistency_policy {
    consistency_level = "Session"
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}
