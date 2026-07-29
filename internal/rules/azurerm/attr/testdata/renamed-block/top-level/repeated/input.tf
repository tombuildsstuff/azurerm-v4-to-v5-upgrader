resource "azurerm_cdn_frontdoor_rule" "example" {
  request_header_action {
    header_action = "Append"
    value         = "one"
  }
  request_header_action {
    header_action = "Overwrite"
    value         = "two"
  }
}
