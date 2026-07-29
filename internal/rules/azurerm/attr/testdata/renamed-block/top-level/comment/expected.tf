resource "azurerm_cdn_frontdoor_rule" "example" {
  # TODO: revisit this header rewrite once the CDN team confirms the new name
  modify_request_header {
    operator     = "Append"
    header_value = "bar"
  }
}
