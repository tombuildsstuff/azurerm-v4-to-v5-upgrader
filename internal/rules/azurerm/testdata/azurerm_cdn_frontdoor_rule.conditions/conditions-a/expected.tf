resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "example"
  cdn_frontdoor_rule_set_id = "/subscriptions/x/ruleSets/rs"
  order                     = 1
  behaviour_on_match        = "Continue"

  conditions {
    client_port {
      operator = "Equal"
      values   = ["8080"]
    }
    request_cookies {
      name     = "session"
      operator = "Equal"
      values   = ["abc"]
    }
    host_name {
      operator = "Equal"
      values   = ["contoso.com"]
    }
    http_version {
      values = ["2.0"]
    }
    device_type {
      values = ["Mobile"]
    }
    post_argument {
      name   = "field"
      values = ["value"]
    }
    query_string {
      operator = "Contains"
      values   = ["foo=bar"]
    }
    remote_address {
      operator = "Any"
      values   = []
    }
    request_body {
      operator = "Contains"
      values   = ["payload"]
    }
    request_header {
      name     = "X-Custom"
      operator = "Equal"
      values   = ["value"]
    }
    request_method {
      values = ["GET"]
    }
    request_scheme {
      values = ["HTTPS"]
    }
    request_url {
      operator = "Contains"
      values   = ["/api"]
    }
    server_port {
      operator = "Equal"
      values   = ["443"]
    }
    socket_address {
      operator = "Any"
      values   = []
    }
    ssl_protocol {
      values = ["TLSv1.2"]
    }
    request_file_extension {
      operator = "Equal"
      values   = ["html"]
    }
    request_filename {
      operator = "Equal"
      values   = ["index"]
    }
    request_path {
      operator = "Contains"
      values   = ["/blog"]
    }
  }
}
