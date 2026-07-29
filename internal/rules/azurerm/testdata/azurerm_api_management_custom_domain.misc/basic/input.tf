resource "azurerm_api_management_custom_domain" "example" {
  api_management_id = azurerm_api_management.example.id

  developer_portal {
    host_name    = "developer.contoso.com"
    key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
  }
  gateway {
    host_name    = "api.contoso.com"
    key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
  }
  management {
    host_name    = "mgmt.contoso.com"
    key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
  }
  portal {
    host_name    = "portal.contoso.com"
    key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
  }
  scm {
    host_name    = "scm.contoso.com"
    key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
  }
}
