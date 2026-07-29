resource "azurerm_virtual_network" "example" {
  name                = "example"
  address_space       = ["10.0.0.0/16"]
  location            = "westeurope"
  resource_group_name = "example"
}

resource "azurerm_subnet" "example" {
  name                 = "example"
  resource_group_name  = "example"
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "managedinstancedelegation"

    service_delegation {
      name = "Microsoft.Sql/managedInstances"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "vm" {
  name                 = "vm"
  resource_group_name  = "example"
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example"
  location            = "westeurope"
  resource_group_name = "example"

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.vm.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_windows_virtual_machine" "example" {
  name                  = "example"
  resource_group_name   = "example"
  location              = "westeurope"
  size                  = "Standard_F2"
  admin_username        = "adminuser"
  admin_password        = "ExamplePassword123!"
  network_interface_ids = [azurerm_network_interface.example.id]

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

resource "azurerm_mssql_server" "example" {
  name                         = "example"
  resource_group_name          = "example"
  location                     = "westeurope"
  version                      = "12.0"
  administrator_login          = "exampleadmin"
  administrator_login_password = "ExamplePassword123!"
  minimum_tls_version          = "1.2"
}

resource "azurerm_mssql_database" "example" {
  name      = "example"
  server_id = azurerm_mssql_server.example.id
  sku_name  = "S0"

  long_term_retention_policy {
    week_of_year              = 5
    immutable_backups_enabled = true
  }

  threat_detection_policy {
    state            = "Enabled"
    storage_endpoint = "https://example.blob.core.windows.net/"
  }
}

resource "azurerm_mssql_database_extended_auditing_policy" "example" {
  database_id                             = azurerm_mssql_database.example.id
  storage_endpoint                        = "https://example.blob.core.windows.net/"
  storage_account_access_key              = "example-access-key"
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 6
}

resource "azurerm_mssql_server_extended_auditing_policy" "example" {
  server_id                  = azurerm_mssql_server.example.id
  storage_endpoint           = "https://example.blob.core.windows.net/"
  storage_account_access_key = "example-access-key"
  retention_in_days          = 6
}

resource "azurerm_mssql_server_security_alert_policy" "example" {
  resource_group_name  = "example"
  server_name          = azurerm_mssql_server.example.name
  state                = "Enabled"
  email_account_admins = true
}

resource "azurerm_mssql_managed_instance" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  license_type        = "LicenseIncluded"
  sku_name            = "GP_Gen5"
  storage_size_in_gb  = 32
  subnet_id           = azurerm_subnet.example.id
  vcores              = 4

  administrator_login          = "exampleadmin"
  administrator_login_password = "ExamplePassword123!"
}

resource "azurerm_mssql_managed_database" "example" {
  name                = "example"
  managed_instance_id = azurerm_mssql_managed_instance.example.id

  long_term_retention_policy {
    week_of_year              = 5
    immutable_backups_enabled = true
  }

  short_term_retention_days = 7
}

resource "azurerm_mssql_virtual_machine" "example" {
  virtual_machine_id = azurerm_windows_virtual_machine.example.id
  sql_license_type   = "PAYG"

  auto_backup {
    retention_period_in_days   = 30
    storage_blob_endpoint      = "https://example.blob.core.windows.net/"
    storage_account_access_key = "ZXhhbXBsZS1hY2Nlc3Mta2V5"
    encryption_enabled         = true
    encryption_password        = "ExamplePassword123!"
  }
}
