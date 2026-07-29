resource "azurerm_kubernetes_cluster_node_pool" "example" {
  name                  = "internal"
  kubernetes_cluster_id = "example-id"
  vm_size               = "Standard_D2_v2"

  kubelet_config {
    container_log_max_files = 8
    cpu_manager_policy      = "static"
  }

  linux_os_config {
    transparent_huge_page = "madvise"
    swap_file_size_mb     = 200
  }
}
