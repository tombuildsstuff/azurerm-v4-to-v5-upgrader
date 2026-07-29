resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = "example"
  location                 = "westeurope"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                 = "example"
  storage_account_name = azurerm_storage_account.example.name
}
