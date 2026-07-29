resource "azurerm_cdn_frontdoor_rule" "example" {
  conditions {
    request_cookies {
      cookie_name = "session"
      values      = ["a", "b"]
    }
  }
}
