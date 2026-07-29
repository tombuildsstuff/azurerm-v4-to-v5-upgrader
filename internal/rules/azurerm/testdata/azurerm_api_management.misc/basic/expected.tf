resource "azurerm_api_management" "example" {
  name = "example"

  hostname_configuration {
    developer_portal {
      host_name                = "developer.contoso.com"
      key_vault_certificate_id = "https://example.vault.azure.net/secrets/cert/abc123"
    }
    management {
      host_name                = "mgmt.contoso.com"
      key_vault_certificate_id = "https://example.vault.azure.net/secrets/cert/abc123"
    }
    portal {
      host_name                = "portal.contoso.com"
      key_vault_certificate_id = "https://example.vault.azure.net/secrets/cert/abc123"
    }
    proxy {
      host_name                = "api.contoso.com"
      key_vault_certificate_id = "https://example.vault.azure.net/secrets/cert/abc123"
    }
    scm {
      host_name                = "scm.contoso.com"
      key_vault_certificate_id = "https://example.vault.azure.net/secrets/cert/abc123"
    }
  }

  protocols {
    http2_enabled = true
  }

  security {
    backend_tls10_enabled  = false
    backend_tls11_enabled  = false
    frontend_ssl30_enabled = false
    frontend_tls10_enabled = false
    frontend_tls11_enabled = false
  }
}
