# azurerm_bot_channels_registration: the parent Bot Service resource that the
# azurerm_bot_channel_ms_teams channel below attaches to. microsoft_app_type
# is included explicitly since it is now Required in v5 (requiredAttrRule
# flags it NeedsReview only when absent; supplying it is fixture
# completeness, not a rule-gap workaround).
resource "azurerm_bot_channels_registration" "example" {
  name                = "example"
  location            = "global"
  resource_group_name = "example"
  sku                 = "F0"
  microsoft_app_id    = "00000000-0000-0000-0000-000000000000"
  microsoft_app_type  = "MultiTenant"
}

# azurerm_bot_channel_ms_teams: exercises the enable_calling -> calling_enabled
# rename (renamedAttrRule).
resource "azurerm_bot_channel_ms_teams" "example" {
  bot_name            = azurerm_bot_channels_registration.example.name
  location            = "global"
  resource_group_name = "example"

  calling_enabled = true
}

# azurerm_bot_service_azure_bot: microsoft_app_type supplied explicitly (see
# note on azurerm_bot_channels_registration above).
resource "azurerm_bot_service_azure_bot" "example" {
  name                = "example"
  location            = "global"
  resource_group_name = "example"
  sku                 = "F0"
  microsoft_app_id    = "00000000-0000-0000-0000-000000000000"
  microsoft_app_type  = "MultiTenant"
}

# azurerm_bot_web_app: microsoft_app_type supplied explicitly (see note on
# azurerm_bot_channels_registration above).
resource "azurerm_bot_web_app" "example" {
  name                = "example-2"
  location            = "global"
  resource_group_name = "example"
  sku                 = "F0"
  microsoft_app_id    = "11111111-1111-1111-1111-111111111111"
  microsoft_app_type  = "MultiTenant"
}

# azurerm_dashboard_grafana: sku pinned to "Standard" and grafana_major_version
# to "12" (both GA-valid) rather than the per-rule seed fixture's
# sku = "Essential" / grafana_major_version = "11"
# (internal/rules/azurerm/testdata/azurerm_dashboard_grafana.misc), which
# valueChangeRule flags as removed values in v5 without rewriting them. Using
# the actually-removed values would make the corpus fail GA validate even
# pre-upgrade (GA v5.0.0 requires sku = "Standard" and grafana_major_version
# in {"12", "13"} - verified against
# internal/services/dashboard/dashboard_grafana_resource.go at the v5.0.0 tag
# in the local hashicorp/terraform-provider-azurerm checkout). Same
# value-removed coverage pattern used for azurerm_synapse_spark_pool in
# modules/dataplatform/main.tf.
resource "azurerm_dashboard_grafana" "example" {
  name                  = "example"
  resource_group_name   = "example"
  location              = "westeurope"
  sku                   = "Standard"
  grafana_major_version = "12"
}

# azurerm_powerbi_embedded: mode deliberately omitted so the upgrader writes
# it explicitly (changedDefaultRule preserves the v4 default of "Gen1" with a
# NOTE comment, since v5 changes the default to "Gen2").
resource "azurerm_powerbi_embedded" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  sku_name            = "A1"
  administrators      = ["admin@example.com"]
  mode                = "Gen1" # NOTE: AzureRM v5 changes this default to 'Gen2'
}

# azurerm_communication_service: data_location supplied explicitly since it is
# now Required in v5 with no default (requiredAttrRule flags it NeedsReview
# only when absent; supplying it is fixture completeness, not a rule-gap
# workaround).
resource "azurerm_communication_service" "example" {
  name                = "example"
  resource_group_name = "example"
  data_location       = "United States"
}

# azurerm_datadog_monitor: minimal sibling required by
# azurerm_datadog_monitor_sso_configuration below.
resource "azurerm_datadog_monitor" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  sku_name            = "Linked"

  datadog_organization {
    api_key         = "00000000000000000000000000000000"
    application_key = "00000000000000000000000000000000"
  }

  user {
    name  = "example"
    email = "admin@example.com"
  }
}

# azurerm_datadog_monitor_sso_configuration: exercises the
# single_sign_on_enabled -> single_sign_on rename (renamedAttrRule).
# enterprise_application_id is required by GA independent of any v5 change
# (unaffected schema field), so it's supplied for fixture completeness.
resource "azurerm_datadog_monitor_sso_configuration" "example" {
  name               = "default"
  datadog_monitor_id = azurerm_datadog_monitor.example.id

  enterprise_application_id = "22222222-2222-2222-2222-222222222222"
  single_sign_on            = "Enable"
}

# azurerm_storage_account: target resource for azurerm_spring_cloud_connection
# below.
resource "azurerm_storage_account" "example" {
  name                            = "exampleappservicesacc"
  resource_group_name             = "example"
  location                        = "westeurope"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true # NOTE: AzureRM v5 changes this default to 'false'
}

# azurerm_spring_cloud_service / azurerm_spring_cloud_app /
# azurerm_spring_cloud_java_deployment: minimal chain required by
# azurerm_spring_cloud_connection below (spring_cloud_id must resolve to a
# Spring Cloud Deployment resource ID). All three are deprecated-but-present
# in v5.0.0 (Azure Spring Apps retirement warning only, not a validate
# error), so they still validate cleanly.
resource "azurerm_spring_cloud_service" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
}

resource "azurerm_spring_cloud_app" "example" {
  name                = "example"
  resource_group_name = "example"
  service_name        = azurerm_spring_cloud_service.example.name
}

resource "azurerm_spring_cloud_java_deployment" "example" {
  name                = "example"
  spring_cloud_app_id = azurerm_spring_cloud_app.example.id
}

# azurerm_spring_cloud_connection: client_type pinned to "dotnet" (GA-valid)
# rather than the per-rule seed fixture's client_type = "none"
# (internal/rules/azurerm/testdata/azurerm_spring_cloud_connection.misc),
# which valueChangeRule flags as a removed value in v5 without rewriting it.
# Using the actually-removed value would make the corpus fail GA validate
# even pre-upgrade. Same value-removed coverage pattern used for
# azurerm_synapse_spark_pool in modules/dataplatform/main.tf.
resource "azurerm_spring_cloud_connection" "example" {
  name               = "example"
  spring_cloud_id    = azurerm_spring_cloud_java_deployment.example.id
  target_resource_id = azurerm_storage_account.example.id
  client_type        = "dotnet"

  authentication {
    type = "systemAssignedIdentity"
  }
}
