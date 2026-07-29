provider "azurerm" {
  features {
  }

  enhanced_validation {
    locations = false
  }

  resource_provider_registrations = "legacy"
}
