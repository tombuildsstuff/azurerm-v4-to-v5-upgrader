resource "azurerm_private_dns_srv_record" "example" {
  name                = "example"
  zone_name           = "example.com"
  resource_group_name = "example"
  ttl                 = 300

  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }
}
