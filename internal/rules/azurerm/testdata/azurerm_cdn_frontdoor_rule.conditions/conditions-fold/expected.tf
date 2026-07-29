resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1
  behaviour_on_match        = "Continue"

  conditions {
    host_name {
      operator = "NotEqual"
      values   = ["contoso.com"]
    }
    query_string {
      values           = ["foo=bar"]
      negate_condition = true
    }
  }
}
