resource "azurerm_site_recovery_replicated_vm" "example" {
  network_interface {
    source_network_interface_id = "/subscriptions/x/nic"
    target_static_ip            = "10.0.0.4"
    target_subnet_name          = "subnet-1"
  }
}
