resource "azurerm_orchestrated_virtual_machine_scale_set" "example" {
  name                        = "example-ovmss"
  resource_group_name         = "example-resources"
  location                    = "West Europe"
  platform_fault_domain_count = 1
  sku_name                    = "Mix"
  instances                   = 1

  sku_profile {
    allocation_strategy = "LowestPrice"
    vm_sizes            = ["Standard_D2s_v3", "Standard_D4s_v3"]
  }

  os_profile {
    windows_configuration {
      admin_username            = "adminuser"
      admin_password            = "P@55w0rd1234!"
      computer_name_prefix      = "vm-"
      automatic_updates_enabled = false
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    caching              = "ReadWrite"
    storage_account_type = "UltraSSD_LRS"
    create_option        = "Empty"
    disk_size_gb         = 10
    lun                  = 0
    disk_iops_read_write = 500
    disk_mbps_read_write = 50
  }

  network_interface {
    name                           = "example"
    primary                        = true
    accelerated_networking_enabled = true
    ip_forwarding_enabled          = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = "/subscriptions/xxx/resourceGroups/example/providers/Microsoft.Network/virtualNetworks/example/subnets/internal"
    }
  }
}
