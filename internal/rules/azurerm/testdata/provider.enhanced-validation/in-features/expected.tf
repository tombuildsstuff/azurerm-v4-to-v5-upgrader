provider "azurerm" {
  features {
    enhanced_validation {
      locations          = true
      resource_providers = true
    }
  }

  resource_provider_registrations = "legacy"
}
