# --- Key Vault ---
# NOTE: the inline `contact` block is intentionally omitted here (not just
# left unset): it was removed in v5 (it used a data-plane API) in favour of
# the separate azurerm_key_vault_certificate_contacts resource, and the
# upgrader only flags it for manual review (flag-only, unrewritten). A v4
# fixture using it would fail GA validate post-upgrade with `Unsupported
# block type: Blocks of type "contact" are not expected here.` Contact
# management is exercised instead via azurerm_key_vault_certificate_contacts
# below. See tests/README.md "Known rule gaps".

resource "azurerm_key_vault" "example" {
  name                      = "example-kv"
  location                  = azurerm_resource_group.example.location
  resource_group_name       = azurerm_resource_group.example.name
  tenant_id                 = "00000000-0000-0000-0000-000000000000"
  sku_name                  = "standard"
  enable_rbac_authorization = true
}

# --- Key Vault certificate contacts ---
# (contact is now a required block in v5; provided explicitly.)

resource "azurerm_key_vault_certificate_contacts" "example" {
  key_vault_id = azurerm_key_vault.example.id

  contact {
    email = "admin@example.com"
    name  = "Admin"
  }
}
