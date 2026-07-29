resource "azurerm_api_management" "example" {
  name = "example"

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
    enable_backend_tls10  = false
    enable_backend_tls11  = false
    enable_frontend_ssl30 = false
    enable_frontend_tls10 = false
    enable_frontend_tls11 = false
  }
}
