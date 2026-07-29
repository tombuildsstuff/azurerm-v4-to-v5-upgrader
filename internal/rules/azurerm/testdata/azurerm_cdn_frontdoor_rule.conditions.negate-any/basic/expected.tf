resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1

  conditions {
    remote_address {
      operator         = "Any"
      negate_condition = true
      values           = ["10.0.0.0/8"]
    }
  }
}
