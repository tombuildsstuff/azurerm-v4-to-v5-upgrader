resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1
  behaviour_on_match        = "Continue"

  conditions {
    remote_address {
      values = ["10.0.0.0/8"]
    }
    socket_address {
      values = ["10.0.0.0/8"]
    }
  }
}
