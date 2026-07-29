resource "azurerm_eventgrid_event_subscription" "eventhub" {
  name        = "example-eventhub-sub"
  scope       = azurerm_resource_group.example.id
  eventhub_id = azurerm_eventhub.example.id
}

resource "azurerm_eventgrid_event_subscription" "hybrid_connection" {
  name                 = "example-hc-sub"
  scope                = azurerm_resource_group.example.id
  hybrid_connection_id = azurerm_relay_hybrid_connection.example.id
}

resource "azurerm_eventgrid_event_subscription" "service_bus_queue" {
  name                 = "example-sbq-sub"
  scope                = azurerm_resource_group.example.id
  service_bus_queue_id = azurerm_servicebus_queue.example.id
}

resource "azurerm_eventgrid_event_subscription" "service_bus_topic" {
  name                 = "example-sbt-sub"
  scope                = azurerm_resource_group.example.id
  service_bus_topic_id = azurerm_servicebus_topic.example.id
}

resource "azurerm_eventgrid_event_subscription" "azure_function" {
  name  = "example-func-sub"
  scope = azurerm_resource_group.example.id

  azure_function_endpoint {
    function_id = azurerm_function_app_function.example.id
  }
}
