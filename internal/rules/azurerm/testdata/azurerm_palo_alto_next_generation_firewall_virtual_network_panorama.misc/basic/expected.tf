resource "azurerm_palo_alto_next_generation_firewall_virtual_network_panorama" "example" {
  name         = "example"
  rulestack_id = "example-rulestack-id"

  network_profile {
    public_ip_address_ids        = ["example-public-ip-id"]
    virtual_hub_id               = "example-virtual-hub-id"
    network_virtual_appliance_id = "example-nva-id"
  }
  plan_id = "panw-cloud-ngfw-payg" # NOTE: AzureRM v5 changes this default to 'panw-cngfw-payg'
}
