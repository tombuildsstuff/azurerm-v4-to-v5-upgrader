resource "azurerm_synapse_spark_pool" "example" {
  name                 = "example"
  synapse_workspace_id = "example-workspace-id"
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"
  spark_version        = "3.2"

  auto_scale {
    max_node_count = 3
    min_node_count = 1
  }
}
