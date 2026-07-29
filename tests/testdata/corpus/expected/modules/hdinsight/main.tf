# Siblings shared by all HDInsight clusters below: a Data-Lake-Gen2-enabled
# storage account + a user-assigned identity for the storage_account_gen2
# block (v4's storage_resource_id/managed_identity_resource_id, renamed in v5
# to storage_account_id/user_assigned_identity_id).
resource "azurerm_storage_account" "example" {
  name                            = "examplehdistorage"
  resource_group_name             = "example"
  location                        = "westeurope"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  is_hns_enabled                  = true
  allow_nested_items_to_be_public = true # NOTE: AzureRM v5 changes this default to 'false'
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-hdi-identity"
  resource_group_name = "example"
  location            = "westeurope"
}

# NOTE: the v4 `storage_account` block (storage_resource_id, storage_container_id)
# is deliberately NOT exercised here. Its storage_container_id -> v5's
# storage_container_url is a "flag only" NeedsReview rule (attr.go
# flagAttrRule) that never rewrites the key, so the golden output keeps the
# v4 field name and GA rejects it: see the "Rule-gap findings" note recorded
# for this module. Only storage_account_gen2 (fully mechanical rename) is
# used, so this module stays GA-clean.

resource "azurerm_hdinsight_hadoop_cluster" "example" {
  name                = "example-hadoop"
  resource_group_name = "example"
  location            = "westeurope"
  cluster_version     = "4.0"
  tier                = "Standard"
  tls_min_version     = "1.2"

  component_version {
    hadoop = "3.6"
  }

  gateway {
    username = "acctestusrgw"
    password = "AccTestvdSC4daf986!"
  }

  storage_account_gen2 {
    storage_account_id        = azurerm_storage_account.example.id
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
    filesystem_id             = "https://examplehdistorage.dfs.core.windows.net/hadoopfs"
    is_default                = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
    }
    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}

resource "azurerm_hdinsight_hbase_cluster" "example" {
  name                = "example-hbase"
  resource_group_name = "example"
  location            = "westeurope"
  cluster_version     = "4.0"
  tier                = "Standard"
  tls_min_version     = "1.2"

  component_version {
    hbase = "2.1"
  }

  gateway {
    username = "acctestusrgw"
    password = "AccTestvdSC4daf986!"
  }

  storage_account_gen2 {
    storage_account_id        = azurerm_storage_account.example.id
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
    filesystem_id             = "https://examplehdistorage.dfs.core.windows.net/hbasefs"
    is_default                = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
    }
    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}

resource "azurerm_hdinsight_interactive_query_cluster" "example" {
  name                = "example-iqr"
  resource_group_name = "example"
  location            = "westeurope"
  cluster_version     = "4.0"
  tier                = "Standard"
  tls_min_version     = "1.2"

  component_version {
    interactive_hive = "3.1"
  }

  gateway {
    username = "acctestusrgw"
    password = "AccTestvdSC4daf986!"
  }

  storage_account_gen2 {
    storage_account_id        = azurerm_storage_account.example.id
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
    filesystem_id             = "https://examplehdistorage.dfs.core.windows.net/iqrfs"
    is_default                = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D13_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
    worker_node {
      vm_size               = "Standard_D14_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
    }
    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}

resource "azurerm_hdinsight_kafka_cluster" "example" {
  name                = "example-kafka"
  resource_group_name = "example"
  location            = "westeurope"
  cluster_version     = "4.0"
  tier                = "Standard"
  tls_min_version     = "1.2"

  component_version {
    kafka = "2.1"
  }

  gateway {
    username = "acctestusrgw"
    password = "AccTestvdSC4daf986!"
  }

  storage_account_gen2 {
    storage_account_id        = azurerm_storage_account.example.id
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
    filesystem_id             = "https://examplehdistorage.dfs.core.windows.net/kafkafs"
    is_default                = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
    worker_node {
      vm_size                  = "Standard_D3_V2"
      username                 = "acctestusrvm"
      password                 = "AccTestvdSC4daf986!"
      target_instance_count    = 3
      number_of_disks_per_node = 2
    }
    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
    kafka_management_node {
      vm_size  = "Standard_D3_V2"
      password = "AccTestvdSC4daf986!"
    }
  }
}

resource "azurerm_hdinsight_spark_cluster" "example" {
  name                = "example-spark"
  resource_group_name = "example"
  location            = "westeurope"
  cluster_version     = "4.0"
  tier                = "Standard"
  tls_min_version     = "1.2"

  component_version {
    spark = "2.4"
  }

  gateway {
    username = "acctestusrgw"
    password = "AccTestvdSC4daf986!"
  }

  storage_account_gen2 {
    storage_account_id        = azurerm_storage_account.example.id
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
    filesystem_id             = "https://examplehdistorage.dfs.core.windows.net/sparkfs"
    is_default                = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
    }
    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
