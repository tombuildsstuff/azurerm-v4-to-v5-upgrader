# --- Log Analytics workspace (azurerm_log_analytics_workspace.internet_access_types.bool-to-enum) ---
# NOTE: local_authentication_disabled is intentionally omitted here (not just
# left unset): the v5 rename to local_authentication_enabled is flag-only
# (upgrader leaves the v4 boolean unchanged and only emits a NeedsReview note
# because the meaning is inverted, not a straight rename). A v4 fixture using
# it would fail GA validate post-upgrade with `Unsupported argument: An
# argument named "local_authentication_disabled" is not expected here.` See
# tests/README.md "Known rule gaps".

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-law"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30

  internet_ingestion_access_type = "Enabled"
  internet_query_access_type     = "Disabled"
}

# --- Log Analytics linked storage account ---

resource "azurerm_storage_account" "example" {
  name                            = "exampleplatformsa"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true # NOTE: AzureRM v5 changes this default to 'false'
}

resource "azurerm_log_analytics_linked_storage_account" "example" {
  data_source_type    = "CustomLogs"
  resource_group_name = azurerm_resource_group.example.name
  workspace_id        = azurerm_log_analytics_workspace.example.id
  storage_account_ids = [azurerm_storage_account.example.id]
}

# --- Application Insights ---
# NOTE: daily_data_cap_notifications_disabled, disable_ip_masking and
# local_authentication_disabled are intentionally omitted (not just left
# unset): all three v5 renames (to *_enabled counterparts, with inverted
# boolean meaning) are flag-only - the upgrader leaves the v4 booleans
# unchanged and only emits NeedsReview notes. A v4 fixture using them would
# fail GA validate post-upgrade with `Unsupported argument`. See
# tests/README.md "Known rule gaps".

resource "azurerm_application_insights" "example" {
  name                = "example-appi"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
  workspace_id        = azurerm_log_analytics_workspace.example.id
}

# --- Monitor diagnostic settings ---
# NOTE: enabled_log.retention_policy (both resources below) and the top-level
# metric block (azurerm_monitor_diagnostic_setting only) are intentionally
# omitted: both were removed in v5 in favour of azurerm_storage_management_policy
# / enabled_metric respectively, and the upgrader only flags them for manual
# review (flag-only, unrewritten). enabled_metric itself is unaffected by any
# v5 rule and already valid pre-5.0, so it is used directly here instead of
# the removed metric block.

resource "azurerm_monitor_diagnostic_setting" "example" {
  name                       = "example-diag"
  target_resource_id         = azurerm_key_vault.example.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id

  enabled_log {
    category = "AuditEvent"
  }

  enabled_metric {
    category = "AllMetrics"
  }
}

resource "azurerm_monitor_aad_diagnostic_setting" "example" {
  name                       = "example-aad-diag"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id

  enabled_log {
    category = "SignInLogs"
  }
}
