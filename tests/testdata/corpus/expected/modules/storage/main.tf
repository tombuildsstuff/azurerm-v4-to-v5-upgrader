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

resource "azurerm_storage_container" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_storage_share" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
  quota              = 50
}

resource "azurerm_storage_table" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_storage_table_entity" "example" {
  storage_table_id = azurerm_storage_table.example.id
  partition_key    = "example-partition"
  row_key          = "example-row"

  entity = {
    example = "example"
  }
}

data "azurerm_storage_container" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id

  depends_on = [azurerm_storage_container.example]
}

data "azurerm_storage_queue" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id

  depends_on = [azurerm_storage_queue.example]
}

data "azurerm_storage_share" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id

  depends_on = [azurerm_storage_share.example]
}

data "azurerm_storage_table" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id

  depends_on = [azurerm_storage_table.example]
}

data "azurerm_storage_table_entity" "example" {
  storage_table_id = azurerm_storage_table.example.id
  partition_key    = "example-partition"
  row_key          = "example-row"

  depends_on = [azurerm_storage_table_entity.example]
}
