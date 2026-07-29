resource "azurerm_linux_virtual_machine_scale_set" "example" {
  name                = "example-vmss"
  resource_group_name = "example-resources"
  location            = "West Europe"
  sku                 = "Standard_D4_v5"
  instances           = 1
  admin_username      = "adminuser"
  upgrade_mode        = "Automatic"

  admin_ssh_key {
    username   = "adminuser"
    public_key = "ssh-rsa AAAAB3NzaC1yc2E example"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  automatic_os_upgrade_policy {
    disable_automatic_rollback  = false
    enable_automatic_os_upgrade = true
  }

  data_disk {
    lun                            = 0
    caching                        = "ReadWrite"
    create_option                  = "Empty"
    disk_size_gb                   = 10
    storage_account_type           = "UltraSSD_LRS"
    ultra_ssd_disk_iops_read_write = 500
    ultra_ssd_disk_mbps_read_write = 50
  }

  network_interface {
    name                          = "example"
    primary                       = true
    enable_accelerated_networking = true
    enable_ip_forwarding          = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = "/subscriptions/xxx/resourceGroups/example/providers/Microsoft.Network/virtualNetworks/example/subnets/internal"
    }
  }
}
