resource "azurerm_cdn_frontdoor_rule" "example" {
  modify_request_header {
    operator     = "Append"
    header_value = "bar"
  }
}
