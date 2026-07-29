# --- API Management ---

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  publisher_name      = "example"
  publisher_email     = "admin@example.com"
  sku_name            = "Developer_1"

  hostname_configuration {
    developer_portal {
      host_name    = "developer.contoso.com"
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
    proxy {
      host_name    = "api.contoso.com"
      key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
    }
    scm {
      host_name    = "scm.contoso.com"
      key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
    }
  }

  protocols {
    enable_http2 = true
  }

  security {
    enable_backend_ssl30  = true
    enable_backend_tls10  = false
    enable_backend_tls11  = false
    enable_frontend_ssl30 = false
    enable_frontend_tls10 = false
    enable_frontend_tls11 = false
  }
}

# --- API Management custom domain ---

resource "azurerm_api_management_custom_domain" "example" {
  api_management_id = azurerm_api_management.example.id

  developer_portal {
    host_name    = "developer2.contoso.com"
    key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
  }
  gateway {
    host_name    = "api2.contoso.com"
    key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
  }
  management {
    host_name    = "mgmt2.contoso.com"
    key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
  }
  portal {
    host_name    = "portal2.contoso.com"
    key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
  }
  scm {
    host_name    = "scm2.contoso.com"
    key_vault_id = "https://example.vault.azure.net/secrets/cert/abc123"
  }
}
