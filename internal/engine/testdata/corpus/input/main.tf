terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
  }
}

# gateway with BGP enabled
resource "azurerm_virtual_network_gateway" "gw" {
  name       = "${var.prefix}-gw"
  type       = "Vpn"
  enable_bgp = true
}

resource "azurerm_kubernetes_cluster" "aks" {
  name                = "${var.prefix}-aks"
  location            = "westeurope"
  resource_group_name = "example"
  dns_prefix          = "example"
  oidc_issuer_enabled = true

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

  node_provisioning_profile {
    mode = "Auto"
  }

  identity {
    type = "SystemAssigned"
  }
}

# storage account exercises a value-change (min_tls_version), a
# changed-default pin (allow_nested_items_to_be_public absent), and a
# removed block (static_website)
resource "azurerm_storage_account" "sa" {
  name                     = "${var.prefix}sa"
  resource_group_name      = "example"
  location                 = "westeurope"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  min_tls_version          = "TLS1_0"

  static_website {
    index_document = "index.html"
  }
}

# name -> id: reference case (auto-migrated)
resource "azurerm_storage_blob" "blob" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.sa.name
  storage_container_name = "content"
  type                   = "Block"
}

# name -> id: literal case (flagged, not auto-migrated)
resource "azurerm_storage_queue" "queue" {
  name                 = "myqueue"
  storage_account_name = "existingsalegacy"
}

resource "azurerm_application_insights" "appi" {
  name                = "${var.prefix}-appi"
  location            = "westeurope"
  resource_group_name = "example"
  application_type    = "web"
  disable_ip_masking  = true
}

# app service was removed in v5 (flagged, not auto-migrated)
resource "azurerm_app_service" "legacy" {
  name = "legacy"
}
