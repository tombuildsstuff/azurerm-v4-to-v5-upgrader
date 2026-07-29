resource "azurerm_storage_table_entity" "example" {
  storage_table_id = azurerm_storage_table.example.id
  partition_key    = "example-partition"
  row_key          = "example-row"

  entity = {
    example = "example"
  }
}
