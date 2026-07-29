resource "azurerm_storage_account" "example" {
  name                            = "example"
  resource_group_name             = "example"
  location                        = "westeurope"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
  min_tls_version                 = "TLS1_2"

  static_website {
    index_document     = "index.html"
    error_404_document = "404.html"
  }
}
