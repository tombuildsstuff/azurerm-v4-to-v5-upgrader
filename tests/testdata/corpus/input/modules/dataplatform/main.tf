# Siblings shared by the Data Factory resources below.
resource "azurerm_data_factory" "example" {
  name                = "example-adf"
  resource_group_name = "example"
  location            = "westeurope"
}

resource "azurerm_storage_account" "example" {
  name                     = "exampledpstorage"
  resource_group_name      = "example"
  location                 = "westeurope"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_data_factory_integration_runtime_self_hosted" "example" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id

  rbac_authorization {
    resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example-adf/integrationRuntimes/example-shared-runtime"
  }
}

# NOTE: key_vault_sas_token block deliberately omitted. The upgrader's
# removedBlockRule for this block only FLAGS it NeedsReview ("key_vault_sas_token
# was renamed to sas_token_linked_key_vault_key in v5") and does not rewrite
# it, so the golden output keeps the v4 block name, which GA rejects
# ("Unsupported block type"). This is the known flag-only-rename category
# (see tests/README.md "Known rule gaps"), not a new gap. connection_string is
# used instead so the resource itself validates GA-clean.
resource "azurerm_data_factory_linked_service_azure_blob_storage" "example" {
  name              = "example"
  data_factory_id   = azurerm_data_factory.example.id
  connection_string = azurerm_storage_account.example.primary_connection_string
}

resource "azurerm_data_factory_linked_service_azure_databricks" "example" {
  name                       = "example"
  data_factory_id            = azurerm_data_factory.example.id
  adb_domain                 = "https://example.azuredatabricks.net"
  msi_work_space_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Databricks/workspaces/example"
  existing_cluster_id        = "1234-567890-abcdefgh"
}

resource "azurerm_data_factory_linked_service_mysql" "example" {
  name              = "example"
  data_factory_id   = azurerm_data_factory.example.id
  connection_string = "Server=example;Port=3306;Database=example;Uid=example;Pwd=example;"
}

resource "azurerm_data_factory_pipeline" "example" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id

  moniter_metrics_after_duration = "PT1H"
}

# Kusto siblings.
# NOTE: language_extensions and virtual_network_configuration blocks are
# deliberately omitted from the cluster below. Both are flag-only
# removedBlockRule findings (language_extensions -> renamed to
# language_extension, not rewritten; virtual_network_configuration -> removed
# entirely in v5, no replacement) that never rewrite the block, so leaving
# either in would make the golden output GA-invalid. Known flag-only-rename /
# flag-only-removal category, not a new gap.
resource "azurerm_kusto_cluster" "example" {
  name                = "examplekusto"
  resource_group_name = "example"
  location            = "westeurope"

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  cluster_name        = azurerm_kusto_cluster.example.name
}

resource "azurerm_kusto_attached_database_configuration" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "westeurope"
  cluster_name        = azurerm_kusto_cluster.example.name
  database_name       = azurerm_kusto_database.example.name
  cluster_resource_id = azurerm_kusto_cluster.example.id
}

# NOTE: azurerm_kusto_eventgrid_data_connection is deliberately EXCLUDED from
# this module. Its `consumer_group` argument was renamed to
# `eventhub_consumer_group_name` in v5 GA, but no rule in
# internal/rules/azurerm/attr/register.go covers this rename at all (only
# eventgrid_resource_id->eventgrid_event_subscription_id and
# managed_identity_resource_id->managed_identity_id are registered). Leaving
# `consumer_group` untouched, the upgraded config fails GA validate with
# "Missing required argument: eventhub_consumer_group_name" +
# "Unsupported argument: consumer_group". This is a genuine RULE GAP
# (missing rename rule), recorded in tests/README.md "Known rule gaps" and
# the task report's "Rule-gap findings".

# Synapse siblings: a Data-Lake-Gen2-enabled storage account + filesystem, and
# a workspace, for the spark pool below.
resource "azurerm_storage_account" "datalake" {
  name                     = "exampledpdatalake"
  resource_group_name      = "example"
  location                 = "westeurope"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = true
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.datalake.id
}

resource "azurerm_synapse_workspace" "example" {
  name                                 = "example-synw"
  resource_group_name                  = "example"
  location                             = "westeurope"
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3-example!"

  identity {
    type = "SystemAssigned"
  }
}

# NOTE: spark_version pinned to "3.4" (GA-valid; GA now only accepts "3.4" or
# "3.5") rather than the per-rule seed fixture's "3.2"
# (internal/rules/azurerm/testdata/azurerm_synapse_spark_pool.misc), which
# valueChangeRule flags as a removed value in v5. This is the value-removed
# coverage pattern: using the actually-removed value would make the corpus
# fail GA validate even pre-upgrade. See task report's "Value-removed
# coverage note". auto_scale.min_node_count is also bumped from the seed
# fixture's 1 to 3 (GA range is 3-200; unrelated to any v5 rule, a plain
# fixture-completeness fix).
resource "azurerm_synapse_spark_pool" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"
  spark_version        = "3.4"

  auto_scale {
    max_node_count = 6
    min_node_count = 3
  }
}
