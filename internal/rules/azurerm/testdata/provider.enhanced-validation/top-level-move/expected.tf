provider "azurerm" {
  features {

    enhanced_validation {
      locations          = false
      resource_providers = true
    }
  }


  resource_provider_registrations = "legacy"
}
