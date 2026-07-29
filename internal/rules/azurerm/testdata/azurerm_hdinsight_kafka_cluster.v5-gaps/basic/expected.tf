resource "azurerm_hdinsight_kafka_cluster" "example" {
  name                = "example"
  resource_group_name = "example-rg"
  location            = "westeurope"
  cluster_version     = "4.0"

  storage_account {
    storage_account_id   = "/subscriptions/x/storageAccounts/sa"
    storage_container_id = "https://sa.blob.core.windows.net/container"
    is_default           = true
  }

  storage_account_gen2 {
    storage_account_id        = "/subscriptions/x/storageAccounts/sa2"
    user_assigned_identity_id = "/subscriptions/x/userAssignedIdentities/id"
    filesystem_id             = "https://sa2.dfs.core.windows.net/fs"
    is_default                = false
  }
}
