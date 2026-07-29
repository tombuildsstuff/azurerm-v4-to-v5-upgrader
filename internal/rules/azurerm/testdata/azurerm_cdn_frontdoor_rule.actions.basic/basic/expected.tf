resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1
  behaviour_on_match        = "Continue"

  actions {
    modify_request_header {
      operator     = "Append"
      header_name  = "X-A"
      header_value = "1"
    }
    modify_response_header {
      operator     = "Overwrite"
      header_name  = "X-B"
      header_value = "2"
    }
    url_redirect {
      redirect_type         = "Found"
      destination_host_name = "contoso.com"
    }
    url_rewrite {
      source_pattern                  = "/a"
      destination_path                = "/b"
      preserve_unmatched_path_enabled = false
    }
  }
}
