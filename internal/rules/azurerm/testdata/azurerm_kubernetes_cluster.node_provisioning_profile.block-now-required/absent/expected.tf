resource "azurerm_kubernetes_cluster" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"
  dns_prefix          = "example"
  oidc_issuer_enabled = true

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
