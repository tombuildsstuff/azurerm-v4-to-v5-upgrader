resource "azurerm_cdn_endpoint_custom_domain" "example" {
  name            = "example"
  cdn_endpoint_id = "example-endpoint-id"
  host_name       = "example.contoso.com"

  cdn_managed_https {
    certificate_type = "Dedicated"
    protocol_type    = "ServerNameIndication"
    tls_version      = "TLS10"
  }
}

resource "azurerm_cdn_endpoint_custom_domain" "example2" {
  name            = "example2"
  cdn_endpoint_id = "example-endpoint-id"
  host_name       = "example2.contoso.com"

  user_managed_https {
    key_vault_secret_id = "https://example.vault.azure.net/secrets/cert/abc123"
    tls_version         = "None"
  }
}
