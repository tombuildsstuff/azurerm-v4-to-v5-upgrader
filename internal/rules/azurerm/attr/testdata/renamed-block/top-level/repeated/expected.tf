resource "azurerm_cdn_frontdoor_rule" "example" {
  modify_request_header {
    operator     = "Append"
    header_value = "one"
  }
  modify_request_header {
    operator     = "Overwrite"
    header_value = "two"
  }
}
