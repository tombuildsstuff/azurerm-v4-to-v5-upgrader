terraform {
  required_providers {
    azurerm = {
      source                = "hashicorp/azurerm"
      version               = "~> 4.0" # keep me
      configuration_aliases = [azurerm.europe]
    }
  }
}
