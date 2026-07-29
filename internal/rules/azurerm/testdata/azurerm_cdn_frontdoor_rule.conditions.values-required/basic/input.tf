resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1

  conditions {
    is_device_condition {
      operator = "Equal"
    }
    remote_address_condition {
      operator = "IPMatch"
    }
  }
}
