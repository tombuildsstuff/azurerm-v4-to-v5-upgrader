# --- Redis cache ---
# (minimum_tls_version uses "1.2", a still-accepted value, rather than the
# removed "1.0"/"1.1" values the azurerm_redis_cache.misc per-rule fixture
# exercises, per the value-removed-category convention. See tests/README.md
# "Known rule gaps".)

resource "azurerm_redis_cache" "example" {
  name                = "example-redis"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  capacity            = 1
  family              = "C"
  sku_name            = "Standard"

  minimum_tls_version = "1.2"
}
