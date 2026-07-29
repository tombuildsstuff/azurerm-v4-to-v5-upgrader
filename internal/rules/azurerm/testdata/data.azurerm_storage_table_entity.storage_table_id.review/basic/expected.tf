data "azurerm_storage_table_entity" "example" {
  storage_table_id = "https://myaccount.table.core.windows.net/mytable"
  partition_key    = "pk"
  row_key          = "rk"
}
