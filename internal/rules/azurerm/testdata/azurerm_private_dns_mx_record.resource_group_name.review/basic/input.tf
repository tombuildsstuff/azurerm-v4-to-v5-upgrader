resource "azurerm_private_dns_mx_record" "example" {
  name                = "example"
  zone_name           = "example.com"
  resource_group_name = "example"
  ttl                 = 300

  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }
}
