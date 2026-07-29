resource "azurerm_data_factory_linked_service_mysql" "example" {
  name              = "example"
  data_factory_id   = "example-data-factory-id"
  connection_string = "Server=example;Port=3306;Database=example;Uid=example;Pwd=example;"
}
