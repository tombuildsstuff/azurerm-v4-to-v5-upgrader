# --- AKS cluster ---
# (oidc_issuer_enabled is set explicitly to preserve v4 behavior, since v5
# changes its default to true. node_provisioning_profile is now a required
# block in v5; provided explicitly.)

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "example-aks"
  oidc_issuer_enabled = true

  default_node_pool {
    name           = "default"
    node_count     = 1
    vm_size        = "Standard_D2_v2"
    vnet_subnet_id = azurerm_subnet.aks.id

    kubelet_config {
      container_log_max_files = 5
      cpu_manager_policy      = "static"
    }

    linux_os_config {
      transparent_huge_page = "always"
      swap_file_size_mb     = 100
    }
  }

  identity {
    type = "SystemAssigned"
  }

  node_provisioning_profile {
    mode = "Auto"
  }
}

# --- AKS node pool ---

resource "azurerm_kubernetes_cluster_node_pool" "example" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id
  vm_size               = "Standard_D2_v2"
  vnet_subnet_id        = azurerm_subnet.aks.id

  kubelet_config {
    container_log_max_files = 8
    cpu_manager_policy      = "static"
  }

  linux_os_config {
    transparent_huge_page = "madvise"
    swap_file_size_mb     = 200
  }
}
