resource "azurerm_data_factory_linked_service_azure_blob_storage" "example" {
  name            = "example"
  data_factory_id = "example-data-factory-id"

  key_vault_sas_token {
    linked_service_name = "example-key-vault-ls"
    secret_name         = "example-sas-token-secret"
  }
}
