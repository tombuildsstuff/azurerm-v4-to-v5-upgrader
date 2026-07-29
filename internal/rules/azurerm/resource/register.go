package resource

import "github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"

func init() {
	rules.Register(removedResourceRule{types: map[string]struct{}{
		"azurerm_ai_services":                                {},
		"azurerm_app_service":                                {},
		"azurerm_app_service_active_slot":                    {},
		"azurerm_app_service_hybrid_connection":              {},
		"azurerm_app_service_plan":                           {},
		"azurerm_app_service_slot":                           {},
		"azurerm_app_service_source_control_token":           {},
		"azurerm_automation_software_update_configuration":   {},
		"azurerm_batch_certificate":                          {},
		"azurerm_data_protection_backup_instance_postgresql": {},
		"azurerm_data_protection_backup_policy_postgresql":   {},
		"azurerm_function_app":                               {},
		"azurerm_function_app_slot":                          {},
		"azurerm_hpc_cache":                                  {},
		"azurerm_hpc_cache_access_policy":                    {},
		"azurerm_hpc_cache_blob_nfs_target":                  {},
		"azurerm_hpc_cache_blob_target":                      {},
		"azurerm_hpc_cache_nfs_target":                       {},
		"azurerm_maps_creator":                               {},
		"azurerm_network_packet_capture":                     {},
		"azurerm_orbital_contact":                            {},
		"azurerm_orbital_contact_profile":                    {},
		"azurerm_orbital_spacecraft":                         {},
		"azurerm_postgresql_active_directory_administrator":  {},
		"azurerm_postgresql_configuration":                   {},
		"azurerm_postgresql_database":                        {},
		"azurerm_postgresql_firewall_rule":                   {},
		"azurerm_postgresql_server":                          {},
		"azurerm_postgresql_server_key":                      {},
		"azurerm_postgresql_virtual_network_rule":            {},
		"azurerm_redis_enterprise_cluster":                   {},
		"azurerm_redis_enterprise_database":                  {},
		"azurerm_restore_point_collection":                   {},
		"azurerm_security_center_auto_provisioning":          {},
		"azurerm_spatial_anchors_account":                    {},
		"azurerm_static_site":                                {},
		"azurerm_static_site_custom_domain":                  {},
	}})

	rules.Register(removedDataSourceRule{types: map[string]struct{}{
		"azurerm_app_service":               {},
		"azurerm_app_service_plan":          {},
		"azurerm_batch_certificate":         {},
		"azurerm_function_app":              {},
		"azurerm_postgresql_server":         {},
		"azurerm_redis_enterprise_database": {},
		"azurerm_spatial_anchors_account":   {},
	}})
}
