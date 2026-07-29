resource "azurerm_site_recovery_replicated_vm" "example" {
  name = "example"

  network_interface {
    source_network_interface_id = "/subscriptions/x/nic"
    ip_configuration {
      target_static_ip = "10.0.0.4"
    }
    ip_configuration {
      target_static_ip = "10.0.0.5"
    }
  }
}
