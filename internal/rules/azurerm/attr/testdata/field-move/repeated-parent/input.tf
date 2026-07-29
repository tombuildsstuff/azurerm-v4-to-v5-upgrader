resource "azurerm_site_recovery_replicated_vm" "example" {
  network_interface {
    source_network_interface_id = "/subscriptions/x/nic1"
    target_static_ip            = "10.0.0.4"
  }

  network_interface {
    source_network_interface_id = "/subscriptions/x/nic2"
    target_static_ip            = "10.0.0.5"
  }
}
