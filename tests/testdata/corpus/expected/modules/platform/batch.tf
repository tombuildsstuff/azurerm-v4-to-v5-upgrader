# --- Batch account + pool ---
# NOTE: the `certificate` block on azurerm_batch_pool is intentionally
# omitted: it was removed in v5 (no longer supported by the Batch service),
# and the upgrader only flags it for manual review (flag-only, unrewritten).
# A v4 fixture using it would fail GA validate post-upgrade with `Unsupported
# block type: Blocks of type "certificate" are not expected here.` See
# tests/README.md "Known rule gaps".

resource "azurerm_batch_account" "example" {
  name                 = "exampleplatformbatch"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_pool" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_batch_account.example.name
  display_name        = "example"
  vm_size             = "Standard_A1"
  node_agent_sku_id   = "batch.node.ubuntu 20.04"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "canonical"
    offer     = "0001-com-ubuntu-server-focal"
    sku       = "20_04-lts"
    version   = "latest"
  }
}
