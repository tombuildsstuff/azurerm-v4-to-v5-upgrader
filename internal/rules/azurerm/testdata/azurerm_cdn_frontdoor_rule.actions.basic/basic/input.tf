resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1
  behavior_on_match         = "Continue"

  actions {
    request_header_action {
      header_action = "Append"
      header_name   = "X-A"
      value         = "1"
    }
    response_header_action {
      header_action = "Overwrite"
      header_name   = "X-B"
      value         = "2"
    }
    url_redirect_action {
      redirect_type        = "Found"
      destination_hostname = "contoso.com"
    }
    url_rewrite_action {
      source_pattern          = "/a"
      destination             = "/b"
      preserve_unmatched_path = false
    }
  }
}
