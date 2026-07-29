resource "azurerm_kubernetes_cluster" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"
  dns_prefix          = "example"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2_v2"

    kubelet_config {
      container_log_max_line = 5
      cpu_manager_policy     = "static"
    }

    linux_os_config {
      transparent_huge_page_enabled = "always"
      swap_file_size_mb             = 100
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
