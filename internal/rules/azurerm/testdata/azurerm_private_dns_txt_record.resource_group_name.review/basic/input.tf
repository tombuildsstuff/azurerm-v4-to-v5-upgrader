resource "azurerm_private_dns_txt_record" "example" {
  name                = "example"
  zone_name           = "example.com"
  resource_group_name = "example"
  ttl                 = 300

  record {
    value = "v=spf1 mx ~all"
  }
}
