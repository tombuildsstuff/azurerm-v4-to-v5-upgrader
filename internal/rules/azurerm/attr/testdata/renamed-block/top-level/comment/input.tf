resource "azurerm_cdn_frontdoor_rule" "example" {
  # TODO: revisit this header rewrite once the CDN team confirms the new name
  request_header_action {
    header_action = "Append"
    value         = "bar"
  }
}
