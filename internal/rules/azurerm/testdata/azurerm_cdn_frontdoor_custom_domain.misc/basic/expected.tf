resource "azurerm_cdn_frontdoor_custom_domain" "example" {
  name                     = "example"
  cdn_frontdoor_profile_id = "example-profile-id"
  host_name                = "example.contoso.com"

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS10"
  }
}
