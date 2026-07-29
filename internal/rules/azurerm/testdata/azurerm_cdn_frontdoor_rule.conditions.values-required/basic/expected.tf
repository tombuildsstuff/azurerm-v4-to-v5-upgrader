resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1

  conditions {
    device_type {
      operator = "Equal"
    }
    remote_address {
      operator = "IPMatch"
    }
  }
}
