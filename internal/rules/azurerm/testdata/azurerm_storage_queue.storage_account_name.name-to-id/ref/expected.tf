resource "azurerm_storage_account" "example" {
  name                            = "example"
  resource_group_name             = "example"
  location                        = "westeurope"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
  min_tls_version                 = "TLS1_2"
}

resource "azurerm_storage_queue" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
}
