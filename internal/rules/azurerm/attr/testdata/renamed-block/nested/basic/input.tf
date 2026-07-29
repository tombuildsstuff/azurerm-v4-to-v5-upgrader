resource "azurerm_cdn_frontdoor_rule" "example" {
  conditions {
    cookies_condition {
      cookie_name      = "session"
      match_values     = ["a", "b"]
      negate_condition = true
    }
  }
}
