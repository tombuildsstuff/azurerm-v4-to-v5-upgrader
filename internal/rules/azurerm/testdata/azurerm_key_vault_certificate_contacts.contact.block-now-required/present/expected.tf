resource "azurerm_key_vault_certificate_contacts" "example" {
  key_vault_id = azurerm_key_vault.example.id

  contact {
    email = "admin@example.com"
    name  = "Admin"
  }
}
