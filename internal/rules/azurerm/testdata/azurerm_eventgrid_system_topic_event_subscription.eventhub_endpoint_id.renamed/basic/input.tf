resource "azurerm_eventgrid_system_topic_event_subscription" "eventhub" {
  name                 = "example-eventhub-sub"
  system_topic         = azurerm_eventgrid_system_topic.example.name
  resource_group_name  = "example-resources"
  eventhub_endpoint_id = azurerm_eventhub.example.id
}

resource "azurerm_eventgrid_system_topic_event_subscription" "hybrid_connection" {
  name                          = "example-hc-sub"
  system_topic                  = azurerm_eventgrid_system_topic.example.name
  resource_group_name           = "example-resources"
  hybrid_connection_endpoint_id = azurerm_relay_hybrid_connection.example.id
}

resource "azurerm_eventgrid_system_topic_event_subscription" "service_bus_queue" {
  name                          = "example-sbq-sub"
  system_topic                  = azurerm_eventgrid_system_topic.example.name
  resource_group_name           = "example-resources"
  service_bus_queue_endpoint_id = azurerm_servicebus_queue.example.id
}

resource "azurerm_eventgrid_system_topic_event_subscription" "service_bus_topic" {
  name                          = "example-sbt-sub"
  system_topic                  = azurerm_eventgrid_system_topic.example.name
  resource_group_name           = "example-resources"
  service_bus_topic_endpoint_id = azurerm_servicebus_topic.example.id
}

resource "azurerm_eventgrid_system_topic_event_subscription" "azure_function" {
  name                = "example-func-sub"
  system_topic        = azurerm_eventgrid_system_topic.example.name
  resource_group_name = "example-resources"

  azure_function_endpoint {
    function_id = azurerm_function_app_function.example.id
  }
}
