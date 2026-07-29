resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1

  actions {
    route_configuration_override_action {
      cdn_frontdoor_origin_group_id = "/subscriptions/x/originGroups/og"
      forwarding_protocol           = "HttpsOnly"
      cache_behavior                = "OverrideAlways"
      cache_duration                = "1.12:00:00"
      compression_enabled           = true
      query_string_caching_behavior = "IgnoreQueryString"
    }
  }
}
