resource "azurerm_user_assigned_identity" "example" {
  name                = "example-identity"
  resource_group_name = "example"
  location            = "westeurope"
}

resource "azurerm_storage_account" "example" {
  name                            = "examplewebappsacc"
  resource_group_name             = "example"
  location                        = "westeurope"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true # NOTE: AzureRM v5 changes this default to 'false'
}

resource "azurerm_service_plan" "linux" {
  name                = "example-linux-plan"
  resource_group_name = "example"
  location            = "westeurope"
  os_type             = "Linux"
  sku_name            = "P1v2"
}

resource "azurerm_service_plan" "windows" {
  name                = "example-windows-plan"
  resource_group_name = "example"
  location            = "westeurope"
  os_type             = "Windows"
  sku_name            = "P1v2"
}

resource "azurerm_linux_web_app" "example" {
  name                = "example-linux-web-app"
  resource_group_name = "example"
  location            = "westeurope"
  service_plan_id     = azurerm_service_plan.linux.id

  site_config {
    application_stack {
      node_version = "18-lts"
    }

    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"
  }
}

resource "azurerm_linux_web_app_slot" "example" {
  name           = "staging"
  app_service_id = azurerm_linux_web_app.example.id

  site_config {
    application_stack {
      node_version = "18-lts"
    }

    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"
  }
}

resource "azurerm_windows_web_app" "example" {
  name                = "example-windows-web-app"
  resource_group_name = "example"
  location            = "westeurope"
  service_plan_id     = azurerm_service_plan.windows.id

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"
  }
}

resource "azurerm_windows_web_app_slot" "example" {
  name           = "staging"
  app_service_id = azurerm_windows_web_app.example.id

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"
  }
}

resource "azurerm_linux_function_app" "example" {
  name                = "example-linux-function-app"
  resource_group_name = "example"
  location            = "westeurope"
  service_plan_id     = azurerm_service_plan.linux.id

  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"
  }
}

resource "azurerm_linux_function_app_slot" "example" {
  name            = "staging"
  function_app_id = azurerm_linux_function_app.example.id

  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"
  }
}

resource "azurerm_windows_function_app" "example" {
  name                = "example-windows-function-app"
  resource_group_name = "example"
  location            = "westeurope"
  service_plan_id     = azurerm_service_plan.windows.id

  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"
  }
}

resource "azurerm_windows_function_app_slot" "example" {
  name            = "staging"
  function_app_id = azurerm_windows_function_app.example.id

  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"
  }
}

resource "azurerm_logic_app_standard" "example" {
  name                       = "example-logic-app"
  resource_group_name        = "example"
  location                   = "westeurope"
  app_service_plan_id        = azurerm_service_plan.windows.id
  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key

  site_config {
    min_tls_version     = "1.2"
    scm_min_tls_version = "1.2"
  }
}

resource "azurerm_static_web_app" "example" {
  name                = "example-static-web-app"
  resource_group_name = "example"
  location            = "westeurope"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }
}
