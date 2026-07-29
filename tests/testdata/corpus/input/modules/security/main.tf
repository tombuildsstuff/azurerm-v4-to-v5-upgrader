resource "azurerm_resource_group" "example" {
  name     = "example-security"
  location = "westeurope"
}

# --- Palo Alto NGFW: Virtual Hub based (local rulestack + panorama) ---

resource "azurerm_virtual_wan" "example" {
  name                = "example-vwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-vhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.0.0/23"
}

resource "azurerm_public_ip" "hub" {
  name                = "example-hub-pip"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_palo_alto_virtual_network_appliance" "example" {
  name           = "example-appliance"
  virtual_hub_id = azurerm_virtual_hub.example.id
}

resource "azurerm_palo_alto_local_rulestack" "example" {
  name                = "example-rulestack"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack" "example" {
  name                = "example-ngfwvh-lr"
  resource_group_name = azurerm_resource_group.example.name
  rulestack_id        = azurerm_palo_alto_local_rulestack.example.id

  network_profile {
    public_ip_address_ids        = [azurerm_public_ip.hub.id]
    virtual_hub_id               = azurerm_virtual_hub.example.id
    network_virtual_appliance_id = azurerm_palo_alto_virtual_network_appliance.example.id
  }
}

resource "azurerm_palo_alto_next_generation_firewall_virtual_hub_panorama" "example" {
  name                   = "example-ngfwvh-pan"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  panorama_base64_config = "VGhpcyBpcyBub3QgYSByZWFsIGNvbmZpZywgcGxlYXNlIHVzZSB5b3VyIFBhbm9yYW1hIHNlcnZlciB0byBnZW5lcmF0ZSBhIHJlYWwgdmFsdWUgZm9yIHRoaXMgcHJvcGVydHkh"

  network_profile {
    public_ip_address_ids        = [azurerm_public_ip.hub.id]
    virtual_hub_id               = azurerm_virtual_hub.example.id
    network_virtual_appliance_id = azurerm_palo_alto_virtual_network_appliance.example.id
  }
}

# --- Palo Alto NGFW: Virtual Network based (local rulestack + panorama) ---

resource "azurerm_virtual_network" "ngfw" {
  name                = "example-ngfw-vnet"
  address_space       = ["10.1.0.0/16"]
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_network_security_group" "ngfw" {
  name                = "example-ngfw-nsg"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_subnet" "trust" {
  name                 = "example-trust-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.ngfw.name
  address_prefixes     = ["10.1.1.0/24"]

  delegation {
    name = "trusted"

    service_delegation {
      name    = "PaloAltoNetworks.Cloudngfw/firewalls"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "trust" {
  subnet_id                 = azurerm_subnet.trust.id
  network_security_group_id = azurerm_network_security_group.ngfw.id
}

resource "azurerm_subnet" "untrust" {
  name                 = "example-untrust-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.ngfw.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "untrusted"

    service_delegation {
      name    = "PaloAltoNetworks.Cloudngfw/firewalls"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "untrust" {
  subnet_id                 = azurerm_subnet.untrust.id
  network_security_group_id = azurerm_network_security_group.ngfw.id
}

resource "azurerm_public_ip" "vnet" {
  name                = "example-vnet-pip"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack" "example" {
  name                = "example-ngfwvn-lr"
  resource_group_name = azurerm_resource_group.example.name
  rulestack_id        = azurerm_palo_alto_local_rulestack.example.id

  network_profile {
    public_ip_address_ids = [azurerm_public_ip.vnet.id]

    vnet_configuration {
      virtual_network_id  = azurerm_virtual_network.ngfw.id
      trusted_subnet_id   = azurerm_subnet.trust.id
      untrusted_subnet_id = azurerm_subnet.untrust.id
    }
  }
}

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_panorama" "example" {
  name                   = "example-ngfwvn-pan"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  panorama_base64_config = "e2RnbmFtZTogY25nZnctYXotZXhhbXBsZSwgdHBsbmFtZTogY25nZnctZXhhbXBsZS10ZW1wbGF0ZS1zdGFjaywgZXhhbXBsZS1wYW5vcmFtYS1zZXJ2ZXI6IDE5Mi4xNjguMC4xLCB2bS1hdXRoLWtleTogMDAwMDAwMDAwMDAwMDAwLCBleHBpcnk6IDIwMjQvMDcvMzF9"

  network_profile {
    public_ip_address_ids = [azurerm_public_ip.vnet.id]

    vnet_configuration {
      virtual_network_id  = azurerm_virtual_network.ngfw.id
      trusted_subnet_id   = azurerm_subnet.trust.id
      untrusted_subnet_id = azurerm_subnet.untrust.id
    }
  }
}

# --- Security Center automation ---

resource "azurerm_logic_app_workflow" "example" {
  name                = "example-logicapp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_security_center_automation" "example" {
  name                = "example-automation"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  scopes              = ["/subscriptions/00000000-0000-0000-0000-000000000000"]

  action {
    type        = "LogicApp"
    resource_id = azurerm_logic_app_workflow.example.id
    trigger_url = "https://example.com/trigger"
  }

  source {
    event_source = "Alerts"
  }
}

# --- Sentinel fusion alert rule ---

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-law"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "example" {
  workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_sentinel_alert_rule_fusion" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
  alert_rule_template_guid   = "00000000-0000-0000-0000-000000000000"

  depends_on = [azurerm_sentinel_log_analytics_workspace_onboarding.example]
}

# --- IoT security solution ---
# NOTE: recommendations_enabled is omitted here (not just left unset) because
# the v5 rename to `recommendations` is flag-only (upgrader leaves the v4
# block name unchanged and only emits a NeedsReview note); a v4 fixture using
# `recommendations_enabled` would fail GA validate post-upgrade with
# `Unsupported block type: Blocks of type "recommendations_enabled" are not
# expected here.` See tests/README.md "Known rule gaps".

resource "azurerm_iothub" "example" {
  name                = "example-iothub-sec"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "S1"
    capacity = 1
  }
}

resource "azurerm_iot_security_solution" "example" {
  name                = "example-iot-security"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  display_name        = "example"
  iothub_ids          = [azurerm_iothub.example.id]
}

# --- Recovery Services vault ---

resource "azurerm_recovery_services_vault" "example" {
  name                = "example-vault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Standard"
  soft_delete_enabled = true
}

# --- Site Recovery: fabrics, containers, policy, mapping ---

resource "azurerm_automation_account" "example" {
  name                = "example-automation-account"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Basic"
}

resource "azurerm_site_recovery_fabric" "primary" {
  name                = "example-fabric-primary"
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_site_recovery_fabric" "secondary" {
  name                = "example-fabric-secondary"
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.example.name
  location            = "northeurope"
}

resource "azurerm_site_recovery_protection_container" "primary" {
  name                 = "example-container-primary"
  resource_group_name  = azurerm_resource_group.example.name
  recovery_vault_name  = azurerm_recovery_services_vault.example.name
  recovery_fabric_name = azurerm_site_recovery_fabric.primary.name
}

resource "azurerm_site_recovery_protection_container" "secondary" {
  name                 = "example-container-secondary"
  resource_group_name  = azurerm_resource_group.example.name
  recovery_vault_name  = azurerm_recovery_services_vault.example.name
  recovery_fabric_name = azurerm_site_recovery_fabric.secondary.name
}

resource "azurerm_site_recovery_replication_policy" "example" {
  name                                                 = "example-policy"
  resource_group_name                                  = azurerm_resource_group.example.name
  recovery_vault_name                                  = azurerm_recovery_services_vault.example.name
  recovery_point_retention_in_minutes                  = 1440
  application_consistent_snapshot_frequency_in_minutes = 240
}

resource "azurerm_site_recovery_protection_container_mapping" "example" {
  name                                      = "example-container-mapping"
  resource_group_name                       = azurerm_resource_group.example.name
  recovery_vault_name                       = azurerm_recovery_services_vault.example.name
  recovery_fabric_name                      = azurerm_site_recovery_fabric.primary.name
  recovery_source_protection_container_name = azurerm_site_recovery_protection_container.primary.name
  recovery_target_protection_container_id   = azurerm_site_recovery_protection_container.secondary.id
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.example.id

  automatic_update {
    enabled               = true
    automation_account_id = azurerm_automation_account.example.id
    authentication_type   = "SystemAssignedIdentity"
  }
}

# --- Site Recovery: replicated VM ---

resource "azurerm_virtual_network" "vm" {
  name                = "example-vm-vnet"
  address_space       = ["10.2.0.0/16"]
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_subnet" "vm_source" {
  name                 = "example-vm-source-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.vm.name
  address_prefixes     = ["10.2.1.0/24"]
}

resource "azurerm_subnet" "vm_target" {
  name                 = "example-vm-target-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.vm.name
  address_prefixes     = ["10.2.2.0/24"]
}

resource "azurerm_public_ip" "vm_source" {
  name                = "example-vm-source-pip"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Basic"
}

resource "azurerm_public_ip" "vm_target" {
  name                = "example-vm-target-pip"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Basic"
}

resource "azurerm_network_interface" "vm" {
  name                = "example-vm-nic"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  ip_configuration {
    name                          = "example"
    subnet_id                     = azurerm_subnet.vm_source.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.vm_source.id
  }
}

resource "azurerm_virtual_machine" "example" {
  name                  = "example-vm"
  resource_group_name   = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  vm_size               = "Standard_B1s"
  network_interface_ids = [azurerm_network_interface.vm.id]

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "example-vm-osdisk"
    os_type           = "Linux"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
  }

  os_profile {
    admin_username = "exampleadmin"
    admin_password = "P@ssw0rd1234!"
    computer_name  = "example-vm"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_storage_account" "staging" {
  name                     = "examplestagingacctsec"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_site_recovery_replicated_vm" "example" {
  name                                      = "example-vm-replication"
  resource_group_name                       = azurerm_resource_group.example.name
  recovery_vault_name                       = azurerm_recovery_services_vault.example.name
  source_recovery_fabric_name               = azurerm_site_recovery_fabric.primary.name
  source_vm_id                              = azurerm_virtual_machine.example.id
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.example.id
  source_recovery_protection_container_name = azurerm_site_recovery_protection_container.primary.name

  target_resource_group_id                = azurerm_resource_group.example.id
  target_recovery_fabric_id               = azurerm_site_recovery_fabric.secondary.id
  target_recovery_protection_container_id = azurerm_site_recovery_protection_container.secondary.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.example.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.staging.id
    target_resource_group_id   = azurerm_resource_group.example.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  network_interface {
    source_network_interface_id   = azurerm_network_interface.vm.id
    target_static_ip              = "10.2.2.4"
    target_subnet_name            = azurerm_subnet.vm_target.name
    recovery_public_ip_address_id = azurerm_public_ip.vm_target.id
  }

  depends_on = [azurerm_site_recovery_protection_container_mapping.example]
}
