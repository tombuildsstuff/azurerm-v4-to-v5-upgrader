resource "azurerm_resource_group" "example" {
  name     = "example-networking"
  location = "westeurope"
}

resource "azurerm_virtual_network" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]

  subnet {
    name             = "internal"
    address_prefixes = ["10.0.1.0/24"]
    service_endpoint {
      service = "Microsoft.Storage"
    }
  }
}

resource "azurerm_subnet" "example" {
  name                 = "example"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  service_endpoint {
    service = "Microsoft.Storage"
  }
  service_endpoint {
    service = "Microsoft.Sql"
  }
}

resource "azurerm_subnet" "gateway" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.3.0/24"]
}

resource "azurerm_public_ip" "gateway" {
  name                = "example-gateway-pip"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  type                = "Vpn"
  vpn_type            = "RouteBased"
  sku                 = "VpnGw1"
  bgp_enabled         = true

  ip_configuration {
    name                          = "example"
    public_ip_address_id          = azurerm_public_ip.gateway.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.gateway.id
  }
}

resource "azurerm_local_network_gateway" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  gateway_address     = "10.0.0.1"
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_virtual_network_gateway_connection" "example" {
  name                       = "example"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.example.id
  local_network_gateway_id   = azurerm_local_network_gateway.example.id
  shared_key                 = "changeme"
  bgp_enabled                = true
}

resource "azurerm_lb" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "example"
    public_ip_address_id = azurerm_public_ip.gateway.id
  }
}

resource "azurerm_lb_backend_address_pool" "example" {
  name            = "example"
  loadbalancer_id = azurerm_lb.example.id
}

resource "azurerm_lb_outbound_rule" "example" {
  name                    = "example"
  loadbalancer_id         = azurerm_lb.example.id
  protocol                = "Tcp"
  backend_address_pool_id = azurerm_lb_backend_address_pool.example.id
  tcp_reset_enabled       = true

  frontend_ip_configuration {
    name = azurerm_lb.example.frontend_ip_configuration[0].name
  }
}

resource "azurerm_private_link_service" "example" {
  name                   = "example"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  proxy_protocol_enabled = true

  nat_ip_configuration {
    name      = "primary"
    primary   = true
    subnet_id = azurerm_subnet.example.id
  }

  load_balancer_frontend_ip_configuration_ids = [
    "${azurerm_lb.example.id}/frontendIPConfigurations/example",
  ]
}

resource "azurerm_network_watcher" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_network_security_group" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_storage_account" "example" {
  name                            = "examplenetworkingsa"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true # NOTE: AzureRM v5 changes this default to 'false'
}

resource "azurerm_network_watcher_flow_log" "example" {
  network_watcher_name = azurerm_network_watcher.example.name
  resource_group_name  = azurerm_resource_group.example.name
  name                 = "example"

  target_resource_id = azurerm_network_security_group.example.id
  storage_account_id = azurerm_storage_account.example.id
  enabled            = true

  retention_policy {
    enabled = true
    days    = 7
  }
}

resource "azurerm_express_route_circuit" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }
}

resource "azurerm_express_route_circuit_peering" "example" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.example.name
  resource_group_name           = azurerm_resource_group.example.name
  vlan_id                       = 100
  primary_peer_address_prefix   = "192.168.10.16/30"
  secondary_peer_address_prefix = "192.168.10.20/30"
  peer_asn                      = 100
}

resource "azurerm_virtual_wan" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.10.0.0/24"
}

resource "azurerm_express_route_gateway" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_hub_id      = azurerm_virtual_hub.example.id
  scale_units         = 1
}

resource "azurerm_express_route_connection" "example" {
  name                             = "example"
  express_route_gateway_id         = azurerm_express_route_gateway.example.id
  express_route_circuit_peering_id = azurerm_express_route_circuit_peering.example.id
  internet_security_enabled        = true
}
