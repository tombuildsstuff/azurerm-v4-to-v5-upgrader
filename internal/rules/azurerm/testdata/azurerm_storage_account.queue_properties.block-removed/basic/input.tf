resource "azurerm_storage_account" "example" {
  name                            = "example"
  resource_group_name             = "example"
  location                        = "westeurope"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = false

  queue_properties {
    cors_rule {
      allowed_origins    = ["http://www.example.com"]
      allowed_methods    = ["GET"]
      allowed_headers    = ["*"]
      exposed_headers    = ["*"]
      max_age_in_seconds = 3600
    }
  }
}
