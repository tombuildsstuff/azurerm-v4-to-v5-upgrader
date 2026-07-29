resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1
  behavior_on_match         = "Continue"

  conditions {
    client_port_condition {
      operator         = "Equal"
      match_values     = ["8080"]
      negate_condition = false
    }
    cookies_condition {
      cookie_name      = "session"
      operator         = "Equal"
      match_values     = ["abc"]
      negate_condition = false
    }
    host_name_condition {
      operator         = "Equal"
      match_values     = ["contoso.com"]
      negate_condition = false
    }
    http_version_condition {
      match_values     = ["2.0"]
      negate_condition = false
    }
    is_device_condition {
      match_values     = ["Mobile"]
      negate_condition = false
    }
    post_args_condition {
      post_args_name   = "field"
      match_values     = ["value"]
      negate_condition = false
    }
    query_string_condition {
      operator         = "Contains"
      match_values     = ["foo=bar"]
      negate_condition = false
    }
    remote_address_condition {
      operator         = "Any"
      match_values     = []
      negate_condition = false
    }
    request_body_condition {
      operator         = "Contains"
      match_values     = ["payload"]
      negate_condition = false
    }
    request_header_condition {
      header_name      = "X-Custom"
      operator         = "Equal"
      match_values     = ["value"]
      negate_condition = false
    }
    request_method_condition {
      match_values     = ["GET"]
      negate_condition = false
    }
    request_scheme_condition {
      match_values     = ["HTTPS"]
      negate_condition = false
    }
    request_uri_condition {
      operator         = "Contains"
      match_values     = ["/api"]
      negate_condition = false
    }
    server_port_condition {
      operator         = "Equal"
      match_values     = ["443"]
      negate_condition = false
    }
    socket_address_condition {
      operator         = "Any"
      match_values     = []
      negate_condition = false
    }
    ssl_protocol_condition {
      match_values     = ["TLSv1.2"]
      negate_condition = false
    }
    url_file_extension_condition {
      operator         = "Equal"
      match_values     = ["html"]
      negate_condition = false
    }
    url_filename_condition {
      operator         = "Equal"
      match_values     = ["index"]
      negate_condition = false
    }
    url_path_condition {
      operator         = "Contains"
      match_values     = ["/blog"]
      negate_condition = false
    }
  }
}
