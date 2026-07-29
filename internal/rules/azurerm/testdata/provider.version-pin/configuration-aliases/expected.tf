terraform {
  required_providers {
    azurerm = {
      source                = "hashicorp/azurerm"
      version               = "=5.0.0" # keep me
      configuration_aliases = [azurerm.europe]
    }
  }
}
