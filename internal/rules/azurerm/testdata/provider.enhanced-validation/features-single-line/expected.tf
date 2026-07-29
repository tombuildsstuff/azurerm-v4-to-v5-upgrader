provider "azurerm" {
  resource_provider_registrations = "legacy"
  features {
    enhanced_validation {
      locations          = false
      resource_providers = true
    }
  }

}
