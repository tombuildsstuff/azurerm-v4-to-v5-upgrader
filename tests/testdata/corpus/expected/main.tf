terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=5.0.0"
    }
  }
}

provider "azurerm" {
  features {}
  resource_provider_registrations = "none"
}

module "appservices" {
  source = "./modules/appservices"
}

module "compute" {
  source = "./modules/compute"
}

module "containers" {
  source = "./modules/containers"
}

module "dataplatform" {
  source = "./modules/dataplatform"
}

module "hdinsight" {
  source = "./modules/hdinsight"
}

module "loadbalancer" {
  source = "./modules/loadbalancer"
}

module "messaging" {
  source = "./modules/messaging"
}

module "mssql" {
  source = "./modules/mssql"
}

module "networking" {
  source = "./modules/networking"
}

module "platform" {
  source = "./modules/platform"
}

module "security" {
  source = "./modules/security"
}

module "storage" {
  source = "./modules/storage"
}

module "webapps" {
  source = "./modules/webapps"
}
