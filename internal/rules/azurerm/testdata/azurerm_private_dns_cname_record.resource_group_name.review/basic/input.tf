resource "azurerm_private_dns_cname_record" "example" {
  name                = "example"
  zone_name           = "example.com"
  resource_group_name = "example"
  ttl                 = 300
  record              = "contoso.com"
}
