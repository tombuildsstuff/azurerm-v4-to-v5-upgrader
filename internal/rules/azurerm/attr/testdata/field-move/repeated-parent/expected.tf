resource "azurerm_site_recovery_replicated_vm" "example" {
  network_interface {
    source_network_interface_id = "/subscriptions/x/nic1"
    ip_configuration {
      target_static_ip = "10.0.0.4"
    }
  }

  network_interface {
    source_network_interface_id = "/subscriptions/x/nic2"
    ip_configuration {
      target_static_ip = "10.0.0.5"
    }
  }
}
