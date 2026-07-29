resource "azurerm_storage_account" "example" {
  name                            = "examplemsgsa"
  resource_group_name             = "example-resources"
  location                        = "westeurope"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
  min_tls_version                 = "TLS1_2"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "example-namespace"
  location            = "westeurope"
  resource_group_name = "example-resources"
  sku                 = "Standard"
  capacity            = 1
  minimum_tls_version = "1.2"
}

resource "azurerm_relay_namespace" "example" {
  name                = "example-relay-namespace"
  location            = "westeurope"
  resource_group_name = "example-resources"
  sku_name            = "Standard"
}

resource "azurerm_relay_hybrid_connection" "example" {
  name                 = "example-hybrid-connection"
  resource_group_name  = "example-resources"
  relay_namespace_name = azurerm_relay_namespace.example.name
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "example-sb-namespace"
  location            = "westeurope"
  resource_group_name = "example-resources"
  sku                 = "Standard"
  minimum_tls_version = "1.2"
}

resource "azurerm_servicebus_queue" "example" {
  name         = "example-queue"
  namespace_id = azurerm_servicebus_namespace.example.id
}

resource "azurerm_servicebus_topic" "example" {
  name         = "example-topic"
  namespace_id = azurerm_servicebus_namespace.example.id
}

resource "azurerm_servicebus_subscription" "example" {
  name               = "example-subscription"
  topic_id           = azurerm_servicebus_topic.example.id
  max_delivery_count = 1
}

# metric_arm_resource_id is intentionally NOT set here: GA v5.0.0 confirms
# it is a Computed-only attribute in the azurerm_eventgrid_system_topic
# schema (`metric_resource_id`, its v5 rename target, is Computed-only too -
# this was never a user-settable field even pre-v5), so setting it produces
# "Value for unconfigurable attribute" regardless of the upgrader's rename.
# This is a fixture-authoring issue, not a rule gap: the
# azurerm_eventgrid_system_topic.metric_arm_resource_id.renamed rule's rename
# itself (metric_arm_resource_id -> metric_resource_id) is mechanically
# correct and remains covered by its per-rule golden fixture under
# internal/rules/azurerm/testdata/azurerm_eventgrid_system_topic.metric_arm_resource_id.renamed.
resource "azurerm_eventgrid_system_topic" "example" {
  name                = "example-system-topic"
  location            = "westeurope"
  resource_group_name = "example-resources"
  source_resource_id  = azurerm_storage_account.example.id
  topic_type          = "Microsoft.Storage.StorageAccounts"
}

resource "azurerm_eventgrid_event_subscription" "eventhub" {
  name        = "example-eventhub-sub"
  scope       = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources"
  eventhub_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.EventHub/namespaces/example-namespace/eventhubs/example-eventhub"
}

resource "azurerm_eventgrid_event_subscription" "hybrid_connection" {
  name                 = "example-hc-sub"
  scope                = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources"
  hybrid_connection_id = azurerm_relay_hybrid_connection.example.id
}

resource "azurerm_eventgrid_event_subscription" "service_bus_queue" {
  name                 = "example-sbq-sub"
  scope                = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources"
  service_bus_queue_id = azurerm_servicebus_queue.example.id
}

resource "azurerm_eventgrid_event_subscription" "service_bus_topic" {
  name                 = "example-sbt-sub"
  scope                = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources"
  service_bus_topic_id = azurerm_servicebus_topic.example.id
}

resource "azurerm_eventgrid_event_subscription" "azure_function" {
  name  = "example-func-sub"
  scope = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources"

  azure_function_endpoint {
    function_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Web/sites/example-func-app/functions/example-function"
  }
}

resource "azurerm_eventgrid_system_topic_event_subscription" "eventhub" {
  name                = "example-eventhub-sub"
  system_topic        = azurerm_eventgrid_system_topic.example.name
  resource_group_name = "example-resources"
  eventhub_id         = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.EventHub/namespaces/example-namespace/eventhubs/example-eventhub"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "hybrid_connection" {
  name                 = "example-hc-sub"
  system_topic         = azurerm_eventgrid_system_topic.example.name
  resource_group_name  = "example-resources"
  hybrid_connection_id = azurerm_relay_hybrid_connection.example.id
}

resource "azurerm_eventgrid_system_topic_event_subscription" "service_bus_queue" {
  name                 = "example-sbq-sub"
  system_topic         = azurerm_eventgrid_system_topic.example.name
  resource_group_name  = "example-resources"
  service_bus_queue_id = azurerm_servicebus_queue.example.id
}

resource "azurerm_eventgrid_system_topic_event_subscription" "service_bus_topic" {
  name                 = "example-sbt-sub"
  system_topic         = azurerm_eventgrid_system_topic.example.name
  resource_group_name  = "example-resources"
  service_bus_topic_id = azurerm_servicebus_topic.example.id
}

resource "azurerm_eventgrid_system_topic_event_subscription" "azure_function" {
  name                = "example-func-sub"
  system_topic        = azurerm_eventgrid_system_topic.example.name
  resource_group_name = "example-resources"

  azure_function_endpoint {
    function_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Web/sites/example-func-app/functions/example-function"
  }
}
