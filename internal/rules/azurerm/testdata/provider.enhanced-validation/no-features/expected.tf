provider "azurerm" {

  resource_provider_registrations = "legacy"
  features {

    enhanced_validation {
      locations          = true
      resource_providers = true
    }
  }
}
