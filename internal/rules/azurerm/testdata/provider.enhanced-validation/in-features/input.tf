provider "azurerm" {
  features {
    enhanced_validation {
    }
  }

  resource_provider_registrations = "legacy"
}
