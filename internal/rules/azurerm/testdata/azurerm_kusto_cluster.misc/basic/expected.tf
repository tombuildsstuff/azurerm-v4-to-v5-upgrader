resource "azurerm_kusto_cluster" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  language_extensions {
    name  = "PYTHON"
    image = "Python3_10_8"
  }

  virtual_network_configuration {
    subnet_id                    = "example-subnet-id"
    engine_public_ip_id          = "example-engine-public-ip-id"
    data_management_public_ip_id = "example-dm-public-ip-id"
  }
}
