resource "azurerm_virtual_network" "example" {
  name                = "example-compute"
  address_space       = ["10.0.0.0/16"]
  location            = "westeurope"
  resource_group_name = "example"
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = "example"
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example-compute-nic"
  location            = "westeurope"
  resource_group_name = "example"

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_key_vault" "example" {
  name                       = "examplecomputekv"
  location                   = "westeurope"
  resource_group_name        = "example"
  tenant_id                  = "00000000-0000-0000-0000-000000000000"
  sku_name                   = "standard"
  rbac_authorization_enabled = true
}

resource "azurerm_key_vault_key" "example" {
  name         = "example-des-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_disk_encryption_set" "example" {
  name                = "example-des"
  resource_group_name = "example"
  location            = "westeurope"
  key_vault_key_id    = azurerm_key_vault_key.example.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_dedicated_host_group" "example" {
  name                        = "example-host-group"
  resource_group_name         = "example"
  location                    = "westeurope"
  platform_fault_domain_count = 2
}

resource "azurerm_dedicated_host" "example" {
  name                    = "example-host"
  location                = "westeurope"
  dedicated_host_group_id = azurerm_dedicated_host_group.example.id
  sku_name                = "DSv3-Type3"
  platform_fault_domain   = 1
  license_type            = "Windows_Server_Hybrid"
}

resource "azurerm_linux_virtual_machine_scale_set" "example" {
  name                = "example-lvmss"
  resource_group_name = "example"
  location            = "westeurope"
  sku                 = "Standard_D4_v5"
  instances           = 1
  admin_username      = "adminuser"
  upgrade_mode        = "Automatic"

  admin_ssh_key {
    username   = "adminuser"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDlOfWzwgtf+JxghdWJOU8DQ2uTUmgiiPgqtgRmez+CluCns2K3BC+r0xdBuzp8Qv1H4D6OVXbI02DJ5kxOuPgxLW7ccw3mlDWBE9r2nJss9zscFEBF8z3z6DqHOUSGvxWgLJfC6qN7LlPHNfpMxhirhW12ACGnIt0YApdUuoVwlq5MwQLh8D5O3JDilNiLQuMVIqXYns5ulpBgqM5ZoNOwCuoXnyva0Ec0Gd5mcXd4Aq8SUJQ02BaGFIybIiOSbVmrtTNWT0Ar4HCb4+zFoHxtWh2s2aG2K2ivcn1nq6m1dpRZaAtAYGrPifmWOCBEPXMkJdKhtbGh7c09UmGNHt79 azurerm-v5-upgrader-corpus"
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

  data_disk {
    lun                  = 0
    caching              = "ReadWrite"
    create_option        = "Empty"
    disk_size_gb         = 10
    storage_account_type = "UltraSSD_LRS"
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
      subnet_id = azurerm_subnet.example.id
    }
  }
}

resource "azurerm_windows_virtual_machine" "example" {
  name                  = "example-wvm"
  resource_group_name   = "example"
  location              = "westeurope"
  size                  = "Standard_F2"
  admin_username        = "adminuser"
  admin_password        = "P@ssw0rd1234!"
  network_interface_ids = [azurerm_network_interface.example.id]

  automatic_updates_enabled = false

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }
}

resource "azurerm_windows_virtual_machine_scale_set" "example" {
  name                 = "example-wvmss"
  resource_group_name  = "example"
  location             = "westeurope"
  sku                  = "Standard_D4_v5"
  instances            = 1
  admin_password       = "P@55w0rd1234!"
  admin_username       = "adminuser"
  computer_name_prefix = "vm-"
  upgrade_mode         = "Automatic"

  automatic_updates_enabled = false

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-Server-Core"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    lun                  = 0
    caching              = "ReadWrite"
    create_option        = "Empty"
    disk_size_gb         = 10
    storage_account_type = "UltraSSD_LRS"
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
      subnet_id = azurerm_subnet.example.id
    }
  }
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "example" {
  name                        = "example-ovmss"
  resource_group_name         = "example"
  location                    = "westeurope"
  platform_fault_domain_count = 1
  sku_name                    = "Mix"
  instances                   = 1

  sku_profile {
    allocation_strategy = "LowestPrice"

    virtual_machine_size {
      name = "Standard_D2s_v3"
    }

    virtual_machine_size {
      name = "Standard_D4s_v3"
    }
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
      subnet_id = azurerm_subnet.example.id
    }
  }
}

resource "azurerm_maintenance_configuration" "example" {
  name                = "example-maintenance"
  resource_group_name = "example"
  location            = "westeurope"
  scope               = "OSImage"
}

resource "azurerm_maintenance_assignment_virtual_machine_scale_set" "example" {
  location                     = "westeurope"
  maintenance_configuration_id = azurerm_maintenance_configuration.example.id
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.example.id
}
