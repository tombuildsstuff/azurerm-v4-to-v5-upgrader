resource "azurerm_redis_cache" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"
  capacity            = 1
  family              = "C"
  sku_name            = "Standard"

  minimum_tls_version = "1.1"
}
