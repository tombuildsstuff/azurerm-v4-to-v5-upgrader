resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1
  behavior_on_match         = "Continue"

  conditions {
    host_name_condition {
      operator         = "Equal"
      match_values     = ["contoso.com"]
      negate_condition = true
    }
    query_string_condition {
      match_values     = ["foo=bar"]
      negate_condition = true
    }
  }
}
