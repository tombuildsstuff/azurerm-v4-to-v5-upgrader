resource "azurerm_site_recovery_replicated_vm" "example" {
  network_interface {
    target_static_ip = "10.0.0.4"
    ip_configuration {
      name = "ipconfig1"
    }
  }
}
