resource "azurerm_windows_virtual_machine" "example" {
  name                  = "example-machine"
  resource_group_name   = "example-resources"
  location              = "West Europe"
  size                  = "Standard_F2"
  admin_username        = "adminuser"
  admin_password        = "P@ssw0rd1234!"
  network_interface_ids = ["/subscriptions/xxx/resourceGroups/example/providers/Microsoft.Network/networkInterfaces/example"]

  enable_automatic_updates = false

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
