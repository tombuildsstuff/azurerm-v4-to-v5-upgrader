resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1

  actions {
    route_configuration_override {
      caching {
        behaviour              = "OverrideAlways"
        duration               = "1.12:00:00"
        compression_enabled    = true
        query_string_behaviour = "IgnoreQueryString"
      }
      origin_group {
        cdn_frontdoor_origin_group_id = "/subscriptions/x/originGroups/og"
        forwarding_protocol           = "HttpsOnly"
      }
    }
  }
}
