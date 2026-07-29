package attr

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"
)

func init() {
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_virtual_network_gateway",
		oldAttr:      "enable_bgp",
		newAttr:      "bgp_enabled",
	})
	rules.Register(removedAttrRule{
		resourceType: "azurerm_recovery_services_vault",
		attr:         "soft_delete_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_api_management",
		blockPath:    []string{"security"},
		oldAttr:      "enable_backend_ssl30",
		newAttr:      "backend_ssl30_enabled",
	})
	storageAccountAllowNestedItemsOldDefault := cty.BoolVal(true)
	rules.Register(changedDefaultRule{
		resourceType: "azurerm_storage_account",
		attr:         "allow_nested_items_to_be_public",
		oldDefault:   &storageAccountAllowNestedItemsOldDefault,
		newDefault:   "false",
	})
	rules.Register(changedDefaultRule{
		resourceType: "azurerm_kubernetes_cluster",
		attr:         "oidc_issuer_enabled",
		oldDefault:   nil,
		newDefault:   "true",
	})
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_storage_account",
		attr:          "min_tls_version",
		removedValues: []string{"TLS1_0", "TLS1_1"},
	})
	rules.Register(nameToIDRule{
		resourceType: "azurerm_storage_container",
		oldAttr:      "storage_account_name",
		newAttr:      "storage_account_id",
		targetType:   "azurerm_storage_account",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_storage_account_customer_managed_key",
		attr:         "key_vault_id",
		reason:       "in v5 the key is specified via key_vault_key_id; consolidate key_name/key_vault_id/key_version into key_vault_key_id",
	})
	rules.Register(removedBlockRule{
		resourceType: "azurerm_storage_account",
		blockPath:    []string{"queue_properties"},
		reason:       "queue_properties was removed from azurerm_storage_account in v5; configure queues via azurerm_storage_account_queue_properties",
	})

	// storage
	rules.Register(flagAttrRule{
		resourceType: "azurerm_storage_account",
		blockPath:    []string{"customer_managed_key"},
		attr:         "managed_hsm_key_id",
		reason:       "in v5 the key is specified via key_vault_key_id; consolidate managed_hsm_key_id into key_vault_key_id",
	})
	rules.Register(removedBlockRule{
		resourceType: "azurerm_storage_account",
		blockPath:    []string{"static_website"},
		reason:       "static_website was removed from azurerm_storage_account in v5; configure via azurerm_storage_account_static_website",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_storage_account_customer_managed_key",
		attr:         "key_name",
		reason:       "in v5 the key is specified via key_vault_key_id; consolidate key_name/key_vault_id/key_version into key_vault_key_id",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_storage_account_customer_managed_key",
		attr:         "key_vault_uri",
		reason:       "in v5 the key is specified via key_vault_key_id; consolidate key_vault_uri/key_name/key_version into key_vault_key_id",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_storage_account_customer_managed_key",
		attr:         "key_version",
		reason:       "in v5 the key is specified via key_vault_key_id; consolidate key_name/key_vault_id/key_version into key_vault_key_id",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_storage_account_customer_managed_key",
		attr:         "managed_hsm_key_id",
		reason:       "in v5 the key is specified via key_vault_key_id; consolidate managed_hsm_key_id into key_vault_key_id",
	})
	rules.Register(nameToIDRule{
		resourceType: "azurerm_storage_blob",
		oldAttr:      "storage_account_name",
		newAttr:      "storage_account_id",
		targetType:   "azurerm_storage_account",
	})
	rules.Register(nameToIDRule{
		resourceType: "azurerm_storage_queue",
		oldAttr:      "storage_account_name",
		newAttr:      "storage_account_id",
		targetType:   "azurerm_storage_account",
	})
	rules.Register(nameToIDRule{
		resourceType: "azurerm_storage_share",
		oldAttr:      "storage_account_name",
		newAttr:      "storage_account_id",
		targetType:   "azurerm_storage_account",
	})
	rules.Register(nameToIDRule{
		resourceType: "azurerm_storage_table",
		oldAttr:      "storage_account_name",
		newAttr:      "storage_account_id",
		targetType:   "azurerm_storage_account",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_storage_share_directory",
		attr:         "storage_share_id",
		reason:       "storage_share_id was replaced by storage_share_url in v5; this is an id-to-url change and must be updated manually",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_storage_share_file",
		attr:         "storage_share_id",
		reason:       "storage_share_id was replaced by storage_share_url in v5; this is an id-to-url change and must be updated manually",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_storage_table_entity",
		attr:         "storage_table_id",
		reason:       "storage_table_id now requires a Resource Manager ID in v5; the deprecated Data Plane URL format is no longer supported and must be updated manually",
	})

	// compute
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_windows_virtual_machine",
		oldAttr:      "enable_automatic_updates",
		newAttr:      "automatic_updates_enabled",
	})
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_dedicated_host",
		attr:          "license_type",
		removedValues: []string{"None"},
	})

	// azurerm_linux_virtual_machine_scale_set
	rules.Register(flagAttrRule{
		resourceType: "azurerm_linux_virtual_machine_scale_set",
		blockPath:    []string{"automatic_os_upgrade_policy"},
		attr:         "disable_automatic_rollback",
		reason:       "disable_automatic_rollback was renamed to automatic_rollback_enabled in v5 and the boolean meaning is inverted; update the value manually (e.g. disable_automatic_rollback = true becomes automatic_rollback_enabled = false)",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_linux_virtual_machine_scale_set",
		blockPath:    []string{"automatic_os_upgrade_policy"},
		oldAttr:      "enable_automatic_os_upgrade",
		newAttr:      "automatic_os_upgrade_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_linux_virtual_machine_scale_set",
		blockPath:    []string{"data_disk"},
		oldAttr:      "ultra_ssd_disk_iops_read_write",
		newAttr:      "disk_iops_read_write",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_linux_virtual_machine_scale_set",
		blockPath:    []string{"data_disk"},
		oldAttr:      "ultra_ssd_disk_mbps_read_write",
		newAttr:      "disk_mbps_read_write",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_linux_virtual_machine_scale_set",
		blockPath:    []string{"network_interface"},
		oldAttr:      "enable_accelerated_networking",
		newAttr:      "accelerated_networking_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_linux_virtual_machine_scale_set",
		blockPath:    []string{"network_interface"},
		oldAttr:      "enable_ip_forwarding",
		newAttr:      "ip_forwarding_enabled",
	})

	// azurerm_windows_virtual_machine_scale_set
	rules.Register(flagAttrRule{
		resourceType: "azurerm_windows_virtual_machine_scale_set",
		blockPath:    []string{"automatic_os_upgrade_policy"},
		attr:         "disable_automatic_rollback",
		reason:       "disable_automatic_rollback was renamed to automatic_rollback_enabled in v5 and the boolean meaning is inverted; update the value manually (e.g. disable_automatic_rollback = true becomes automatic_rollback_enabled = false)",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_windows_virtual_machine_scale_set",
		blockPath:    []string{"automatic_os_upgrade_policy"},
		oldAttr:      "enable_automatic_os_upgrade",
		newAttr:      "automatic_os_upgrade_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_windows_virtual_machine_scale_set",
		blockPath:    []string{"data_disk"},
		oldAttr:      "ultra_ssd_disk_iops_read_write",
		newAttr:      "disk_iops_read_write",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_windows_virtual_machine_scale_set",
		blockPath:    []string{"data_disk"},
		oldAttr:      "ultra_ssd_disk_mbps_read_write",
		newAttr:      "disk_mbps_read_write",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_windows_virtual_machine_scale_set",
		oldAttr:      "enable_automatic_updates",
		newAttr:      "automatic_updates_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_windows_virtual_machine_scale_set",
		blockPath:    []string{"network_interface"},
		oldAttr:      "enable_accelerated_networking",
		newAttr:      "accelerated_networking_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_windows_virtual_machine_scale_set",
		blockPath:    []string{"network_interface"},
		oldAttr:      "enable_ip_forwarding",
		newAttr:      "ip_forwarding_enabled",
	})

	// azurerm_orchestrated_virtual_machine_scale_set
	rules.Register(flagAttrRule{
		resourceType: "azurerm_orchestrated_virtual_machine_scale_set",
		blockPath:    []string{"sku_profile"},
		attr:         "vm_sizes",
		reason:       "sku_profile.vm_sizes (a list of VM size strings) was removed in v5 in favour of one or more sku_profile.virtual_machine_size blocks; restructure manually",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_orchestrated_virtual_machine_scale_set",
		blockPath:    []string{"os_profile", "windows_configuration"},
		oldAttr:      "enable_automatic_updates",
		newAttr:      "automatic_updates_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_orchestrated_virtual_machine_scale_set",
		blockPath:    []string{"network_interface"},
		oldAttr:      "enable_accelerated_networking",
		newAttr:      "accelerated_networking_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_orchestrated_virtual_machine_scale_set",
		blockPath:    []string{"network_interface"},
		oldAttr:      "enable_ip_forwarding",
		newAttr:      "ip_forwarding_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_orchestrated_virtual_machine_scale_set",
		blockPath:    []string{"data_disk"},
		oldAttr:      "ultra_ssd_disk_iops_read_write",
		newAttr:      "disk_iops_read_write",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_orchestrated_virtual_machine_scale_set",
		blockPath:    []string{"data_disk"},
		oldAttr:      "ultra_ssd_disk_mbps_read_write",
		newAttr:      "disk_mbps_read_write",
	})

	// networking
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_virtual_network_gateway_connection",
		oldAttr:      "enable_bgp",
		newAttr:      "bgp_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_lb_nat_rule",
		oldAttr:      "enable_floating_ip",
		newAttr:      "floating_ip_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_lb_nat_rule",
		oldAttr:      "enable_tcp_reset",
		newAttr:      "tcp_reset_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_lb_outbound_rule",
		oldAttr:      "enable_tcp_reset",
		newAttr:      "tcp_reset_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_application_gateway",
		oldAttr:      "enable_http2",
		newAttr:      "http2_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_application_gateway",
		blockPath:    []string{"ssl_profile"},
		oldAttr:      "verify_client_cert_issuer_dn",
		newAttr:      "verify_client_certificate_issuer_dn",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_express_route_connection",
		oldAttr:      "enable_internet_security",
		newAttr:      "internet_security_enabled",
	})
	rules.Register(removedAttrRule{
		resourceType: "azurerm_express_route_connection",
		attr:         "private_link_fast_path_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_private_link_service",
		oldAttr:      "enable_proxy_protocol",
		newAttr:      "proxy_protocol_enabled",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_network_watcher_flow_log",
		oldAttr:      "network_security_group_id",
		newAttr:      "target_resource_id",
	})

	// azurerm_private_dns_*_record: v5 replaces resource_group_name + zone_name
	// with a single private_dns_zone_id. Both old attributes must be
	// consolidated by hand into that one id (there's no mechanical way to
	// derive the zone's id from a name/resource-group pair alone), so each is
	// flagged individually rather than auto-renamed.
	for _, resourceType := range []string{
		"azurerm_private_dns_a_record",
		"azurerm_private_dns_aaaa_record",
		"azurerm_private_dns_cname_record",
		"azurerm_private_dns_mx_record",
		"azurerm_private_dns_ptr_record",
		"azurerm_private_dns_srv_record",
		"azurerm_private_dns_txt_record",
	} {
		rules.Register(flagAttrRule{
			resourceType: resourceType,
			attr:         "resource_group_name",
			reason:       "resource_group_name and zone_name were removed in v5 in favour of private_dns_zone_id; consolidate both into private_dns_zone_id",
		})
		rules.Register(flagAttrRule{
			resourceType: resourceType,
			attr:         "zone_name",
			reason:       "zone_name and resource_group_name were removed in v5 in favour of private_dns_zone_id; consolidate both into private_dns_zone_id",
		})
	}
	rules.Register(flagAttrRule{
		resourceType: "azurerm_private_dns_zone_virtual_network_link",
		attr:         "private_dns_zone_name",
		reason:       "private_dns_zone_name and resource_group_name were removed in v5 in favour of private_dns_zone_id; consolidate both into private_dns_zone_id",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_private_dns_zone_virtual_network_link",
		attr:         "resource_group_name",
		reason:       "resource_group_name and private_dns_zone_name were removed in v5 in favour of private_dns_zone_id; consolidate both into private_dns_zone_id",
	})

	// kubernetes
	//
	// azurerm_kubernetes_cluster.oidc_issuer_enabled (changed-default) is
	// already registered above; not duplicated here.
	//
	// azurerm_kubernetes_cluster.node_provisioning_profile becomes a
	// required block in v5 (it's Optional+Computed pre-5.0, Required in
	// 5.0). The tool can't supply a value, so this flags the block as
	// needing manual review when it's absent.
	rules.Register(requiredBlockRule{
		resourceType: "azurerm_kubernetes_cluster",
		blockPath:    []string{"node_provisioning_profile"},
		reason:       "node_provisioning_profile is required in v5; add it explicitly",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_kubernetes_cluster",
		blockPath:    []string{"default_node_pool", "kubelet_config"},
		oldAttr:      "container_log_max_line",
		newAttr:      "container_log_max_files",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_kubernetes_cluster",
		blockPath:    []string{"default_node_pool", "linux_os_config"},
		oldAttr:      "transparent_huge_page_enabled",
		newAttr:      "transparent_huge_page",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_kubernetes_cluster_node_pool",
		blockPath:    []string{"kubelet_config"},
		oldAttr:      "container_log_max_line",
		newAttr:      "container_log_max_files",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_kubernetes_cluster_node_pool",
		blockPath:    []string{"linux_os_config"},
		oldAttr:      "transparent_huge_page_enabled",
		newAttr:      "transparent_huge_page",
	})

	// keyvault
	//
	// azurerm_key_vault: enable_rbac_authorization was renamed to
	// rbac_authorization_enabled (value-preserving bool rename). Separately,
	// the guide also lists rbac_authorization_enabled as now Required in v5
	// (it was Optional+Computed pre-5.0); both rules are registered and
	// coexist deliberately - the rename fires when the old name is used, and
	// the required-attr check flags any resource where, after the rename (or
	// because it was never set at all), rbac_authorization_enabled still
	// ends up absent.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_key_vault",
		oldAttr:      "enable_rbac_authorization",
		newAttr:      "rbac_authorization_enabled",
	})
	rules.Register(requiredAttrRule{
		resourceType: "azurerm_key_vault",
		attr:         "rbac_authorization_enabled",
		reason:       "rbac_authorization_enabled is required in v5; add it explicitly",
	})
	rules.Register(removedBlockRule{
		resourceType: "azurerm_key_vault",
		blockPath:    []string{"contact"},
		reason:       "the deprecated contact block was removed in v5 (it used a data-plane API); manage contacts via the azurerm_key_vault_certificate_contacts resource instead",
	})
	rules.Register(requiredBlockRule{
		resourceType: "azurerm_key_vault_certificate_contacts",
		blockPath:    []string{"contact"},
		reason:       "contact is required in v5; add at least one contact block explicitly",
	})

	// azurerm_federated_identity_credential: parent_id was removed in favour
	// of user_assigned_identity_id. Both attributes hold the same User
	// Assigned Identity resource ID value (parent_id was a deprecated alias,
	// not a name), so this is a straightforward value-preserving id-to-id
	// rename.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_federated_identity_credential",
		oldAttr:      "parent_id",
		newAttr:      "user_assigned_identity_id",
	})

	// azurerm_security_center_automation: action.type no longer accepts the
	// legacy lowercase values; the tool can't safely pick the replacement on
	// the caller's behalf, so this only flags the value for manual review.
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_security_center_automation",
		blockPath:     []string{"action"},
		attr:          "type",
		removedValues: []string{"logicapp", "eventhub", "loganalytics"},
	})

	// azurerm_disk_encryption_set: the deprecated managed_hsm_key_id
	// property was removed in favour of key_vault_key_id. Although both
	// attributes hold the same shape of ID value and are mutually exclusive
	// (ExactlyOneOf) pre-5.0, this is treated as a consolidation rather than
	// a mechanical rename - consistent with how the analogous
	// azurerm_storage_account_customer_managed_key.managed_hsm_key_id case
	// above is handled - so it's flagged for manual review rather than
	// auto-renamed.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_disk_encryption_set",
		attr:         "managed_hsm_key_id",
		reason:       "in v5 the key is specified via key_vault_key_id; consolidate managed_hsm_key_id into key_vault_key_id",
	})

	// app
	//
	// remote_debugging_version no longer accepts VS2017/VS2019 as a value on
	// any Linux/Windows Function App or Web App resource (or their slots);
	// the tool can't safely pick a replacement, so this only flags the value
	// for manual review.
	for _, resourceType := range []string{
		"azurerm_linux_function_app",
		"azurerm_linux_function_app_slot",
		"azurerm_windows_function_app",
		"azurerm_windows_function_app_slot",
		"azurerm_linux_web_app",
		"azurerm_linux_web_app_slot",
		"azurerm_windows_web_app",
		"azurerm_windows_web_app_slot",
	} {
		rules.Register(valueChangeRule{
			resourceType:  resourceType,
			blockPath:     []string{"site_config"},
			attr:          "remote_debugging_version",
			removedValues: []string{"VS2017", "VS2019"},
		})
	}

	// azurerm_linux_web_app / azurerm_linux_web_app_slot: the deprecated
	// site_config.application_stack.ruby_version property was removed in v5
	// (Ruby is no longer offered as a Linux Web App runtime).
	for _, resourceType := range []string{
		"azurerm_linux_web_app",
		"azurerm_linux_web_app_slot",
	} {
		rules.Register(removedAttrRule{
			resourceType: resourceType,
			blockPath:    []string{"site_config", "application_stack"},
			attr:         "ruby_version",
		})
	}

	// azurerm_windows_web_app / azurerm_windows_web_app_slot:
	// virtual_network_image_pull_enabled is Optional+Computed pre-5.0 (its
	// v4 default value tracks the API/environment rather than a fixed
	// literal, so it isn't confidently known) and becomes Optional with a
	// fixed default of false in v5.
	for _, resourceType := range []string{
		"azurerm_windows_web_app",
		"azurerm_windows_web_app_slot",
	} {
		rules.Register(changedDefaultRule{
			resourceType: resourceType,
			attr:         "virtual_network_image_pull_enabled",
			oldDefault:   nil,
			newDefault:   "false",
		})
	}

	// azurerm_logic_app_standard
	//
	// app_service_plan_id now validates that its value is an App Service
	// Plan ID, and that validation is case-sensitive; the tool can't verify
	// casing on the caller's behalf, so this is flagged for manual review
	// rather than auto-migrated.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_logic_app_standard",
		attr:         "app_service_plan_id",
		reason:       "app_service_plan_id now validates its value as an App Service Plan ID (case-sensitive) in v5; confirm the value is a correctly-cased Resource Manager ID",
	})
	// client_certificate_mode now defaults to Required (matching the API
	// default) instead of being left unset pre-5.0; the v4 default isn't a
	// fixed literal, so it isn't confidently known.
	rules.Register(changedDefaultRule{
		resourceType: "azurerm_logic_app_standard",
		attr:         "client_certificate_mode",
		oldDefault:   nil,
		newDefault:   "Required",
	})
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_logic_app_standard",
		blockPath:     []string{"site_config"},
		attr:          "min_tls_version",
		removedValues: []string{"1.0", "1.1"},
	})
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_logic_app_standard",
		blockPath:     []string{"site_config"},
		attr:          "scm_min_tls_version",
		removedValues: []string{"1.0", "1.1"},
	})
	// site_config.public_network_access_enabled (bool) was removed in v5 in
	// favour of the top-level public_network_access (string:
	// Enabled/Disabled) attribute - both a type change (bool -> string) and
	// a location change (nested -> top-level), so this isn't a mechanical
	// rename and is flagged for manual review instead.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_logic_app_standard",
		blockPath:    []string{"site_config"},
		attr:         "public_network_access_enabled",
		reason:       "the deprecated site_config.public_network_access_enabled (bool) was removed in v5 in favour of the top-level public_network_access (string: Enabled/Disabled); this is a type and location change, migrate manually",
	})

	// messaging
	//
	// azurerm_servicebus_namespace / azurerm_eventhub_namespace:
	// minimum_tls_version no longer accepts "1.0" or "1.1" as a value in v5
	// (only "1.2" remains); the tool can't safely pick a replacement, so this
	// only flags the value for manual review.
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_servicebus_namespace",
		attr:          "minimum_tls_version",
		removedValues: []string{"1.0", "1.1"},
	})
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_eventhub_namespace",
		attr:          "minimum_tls_version",
		removedValues: []string{"1.0", "1.1"},
	})

	// azurerm_eventhub: the deprecated namespace_name + resource_group_name
	// properties were removed in v5 in favour of a single namespace_id
	// property. Both old attributes must be consolidated by hand into that
	// one id (there's no mechanical way to derive the namespace's id from a
	// name/resource-group pair alone), so each is flagged individually
	// rather than auto-renamed.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_eventhub",
		attr:         "namespace_name",
		reason:       "namespace_name and resource_group_name were removed in v5 in favour of namespace_id; consolidate both into namespace_id",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_eventhub",
		attr:         "resource_group_name",
		reason:       "resource_group_name and namespace_name were removed in v5 in favour of namespace_id; consolidate both into namespace_id",
	})

	// azurerm_eventgrid_event_subscription / azurerm_eventgrid_system_topic_event_subscription:
	// the deprecated *_endpoint_id properties were removed in v5 in favour of
	// their *_id equivalents. Both the old and new attribute names are backed
	// by the same schema helpers upstream (same shape, same ID validation),
	// so this is a straightforward value-preserving id-to-id rename.
	for _, resourceType := range []string{
		"azurerm_eventgrid_event_subscription",
		"azurerm_eventgrid_system_topic_event_subscription",
	} {
		rules.Register(renamedAttrRule{
			resourceType: resourceType,
			oldAttr:      "eventhub_endpoint_id",
			newAttr:      "eventhub_id",
		})
		rules.Register(renamedAttrRule{
			resourceType: resourceType,
			oldAttr:      "hybrid_connection_endpoint_id",
			newAttr:      "hybrid_connection_id",
		})
		rules.Register(renamedAttrRule{
			resourceType: resourceType,
			oldAttr:      "service_bus_queue_endpoint_id",
			newAttr:      "service_bus_queue_id",
		})
		rules.Register(renamedAttrRule{
			resourceType: resourceType,
			oldAttr:      "service_bus_topic_endpoint_id",
			newAttr:      "service_bus_topic_id",
		})
		// Validation for azure_function_endpoint.function_id now requires an
		// Azure Function resource ID and is case-sensitive; the tool can't
		// verify casing/shape on the caller's behalf, so this is flagged for
		// manual review rather than auto-migrated.
		rules.Register(flagAttrRule{
			resourceType: resourceType,
			blockPath:    []string{"azure_function_endpoint"},
			attr:         "function_id",
			reason:       "azure_function_endpoint.function_id now validates its value as an Azure Function resource ID (case-sensitive) in v5; confirm the value is a correctly-cased Function resource ID rather than a generic Azure resource ID",
		})
	}

	// azurerm_eventgrid_system_topic: the deprecated metric_arm_resource_id
	// and source_arm_resource_id properties were removed in v5 in favour of
	// metric_resource_id and source_resource_id respectively. Both pairs hold
	// the same shape of ID value, so these are value-preserving id-to-id
	// renames. (Note: metric_resource_id/metric_arm_resource_id are
	// Computed-only in the provider schema on both sides of the rename, so
	// they're never actually set as a literal attribute in a resource block
	// in practice; the rule is still registered since it's a documented v5
	// change and is harmless - a no-op - when the attribute isn't present.)
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_eventgrid_system_topic",
		oldAttr:      "metric_arm_resource_id",
		newAttr:      "metric_resource_id",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_eventgrid_system_topic",
		oldAttr:      "source_arm_resource_id",
		newAttr:      "source_resource_id",
	})

	// NOTE: the upgrade guide also lists "The deprecated source_arm_resource_id
	// property has been removed in favour of the source_resource_id property"
	// under azurerm_eventgrid_system_topic_event_subscription. This appears to
	// be a documentation copy/paste error from the azurerm_eventgrid_system_topic
	// section above it: neither source_arm_resource_id nor source_resource_id
	// exist anywhere in the eventgrid_system_topic_event_subscription resource
	// schema in either the v4.50.0 tag or the current main branch of
	// terraform-provider-azurerm. Skipped as unverifiable against source.

	// database
	//
	// azurerm_cosmosdb_account: the deprecated local_authentication_disabled
	// property was removed in v5 in favour of local_authentication_enabled.
	// This is NOT a value-preserving rename: verified against source
	// (internal/services/cosmos/cosmosdb_account_resource.go), pre-5.0
	// `DisableLocalAuth = local_authentication_disabled` directly, while v5
	// sets `DisableLocalAuth = !local_authentication_enabled` - the boolean
	// meaning is inverted (disabled=true means the same thing as
	// enabled=false), so this is flagged for manual review rather than
	// auto-renamed.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_cosmosdb_account",
		attr:         "local_authentication_disabled",
		reason:       "local_authentication_disabled was removed in v5 in favour of local_authentication_enabled, and the boolean meaning is inverted; update the value manually (e.g. local_authentication_disabled = true becomes local_authentication_enabled = false)",
	})
	// The deprecated managed_hsm_key_id property was removed in favour of
	// key_vault_key_id - consolidation precedent, same as
	// azurerm_disk_encryption_set / azurerm_storage_account above.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_cosmosdb_account",
		attr:         "managed_hsm_key_id",
		reason:       "in v5 the key is specified via key_vault_key_id; consolidate managed_hsm_key_id into key_vault_key_id",
	})
	// minimal_tls_version no longer accepts Tls or Tls11 as a value.
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_cosmosdb_account",
		attr:          "minimal_tls_version",
		removedValues: []string{"Tls", "Tls11"},
	})

	// azurerm_mssql_database
	//
	// NOTE: "The enclave_type property is no longer computed" is a
	// computed-behavior change only (no config-visible attribute rename,
	// removal, or value change), so there's nothing for the tool to edit or
	// flag; skipped.
	//
	// long_term_retention_policy.{weekly,monthly,yearly}_retention now
	// default to PT0S in v5. Verified against source
	// (internal/services/mssql/mssql_database_resource.go and
	// helper/sql_retention_policies.go): pre-5.0 these are Optional+Computed
	// with no fixed literal default (the value tracks whatever Azure last
	// returned), so the v4 default isn't confidently known.
	for _, retentionAttr := range []string{"weekly_retention", "monthly_retention", "yearly_retention"} {
		rules.Register(changedDefaultRule{
			resourceType: "azurerm_mssql_database",
			blockPath:    []string{"long_term_retention_policy"},
			attr:         retentionAttr,
			oldDefault:   nil,
			newDefault:   "PT0S",
		})
	}
	// threat_detection_policy.email_account_admins was removed in favour of
	// threat_detection_policy.email_account_admins_enabled. Verified against
	// source: unlike azurerm_mssql_server_security_alert_policy below, this
	// is NOT a value-preserving bool rename - pre-5.0
	// threat_detection_policy.email_account_admins is a TypeString enum
	// ("Enabled"/"Disabled"), while email_account_admins_enabled is a
	// TypeBool. Blindly renaming the key would leave a string literal (e.g.
	// "Enabled") assigned to a bool-typed attribute, which fails type
	// conversion; flagged for manual review instead.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_mssql_database",
		blockPath:    []string{"threat_detection_policy"},
		attr:         "email_account_admins",
		reason:       "email_account_admins was removed in v5 in favour of email_account_admins_enabled; this also changes the value's type from a string (Enabled/Disabled) to a bool, so migrate the value manually",
	})

	// The deprecated long_term_retention_policy.immutable_backups_enabled
	// property was removed in v5 (same field, same reasoning, as
	// azurerm_mssql_managed_database below - the upgrade guide only
	// documents it under the managed_database heading, but the provider
	// source shows it was also present on the plain mssql_database
	// resource pre-5.0 and is likewise absent from the v5 base schema).
	rules.Register(removedAttrRule{
		resourceType: "azurerm_mssql_database",
		blockPath:    []string{"long_term_retention_policy"},
		attr:         "immutable_backups_enabled",
	})

	// azurerm_mssql_database_extended_auditing_policy: the deprecated
	// storage_endpoint property was removed in favour of blob_storage_endpoint.
	// Verified against source: both are plain TypeString URLs holding the
	// same value, so this is a straightforward value-preserving rename.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_mssql_database_extended_auditing_policy",
		oldAttr:      "storage_endpoint",
		newAttr:      "blob_storage_endpoint",
	})

	// azurerm_mssql_elasticpool: the only documented v5 change is
	// "The enclave_type property is no longer computed", a computed-behavior
	// change only; skipped (same reasoning as azurerm_mssql_database above).

	// azurerm_mssql_managed_database
	//
	// The deprecated long_term_retention_policy.immutable_backups_enabled
	// property was removed in v5 (the guide/provider note it was
	// non-functional and mistakenly exposed).
	rules.Register(removedAttrRule{
		resourceType: "azurerm_mssql_managed_database",
		blockPath:    []string{"long_term_retention_policy"},
		attr:         "immutable_backups_enabled",
	})
	// long_term_retention_policy.{weekly,monthly,yearly}_retention now
	// default to PT0S in v5, same Optional+Computed-with-no-fixed-default
	// situation pre-5.0 as azurerm_mssql_database above.
	for _, retentionAttr := range []string{"weekly_retention", "monthly_retention", "yearly_retention"} {
		rules.Register(changedDefaultRule{
			resourceType: "azurerm_mssql_managed_database",
			blockPath:    []string{"long_term_retention_policy"},
			attr:         retentionAttr,
			oldDefault:   nil,
			newDefault:   "PT0S",
		})
	}

	// azurerm_mssql_managed_instance
	//
	// minimum_tls_version no longer accepts 1.0 or 1.1 as a value.
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_mssql_managed_instance",
		attr:          "minimum_tls_version",
		removedValues: []string{"1.0", "1.1"},
	})
	// proxy_override now defaults to Redirect (verified against source:
	// pre-5.0 it's Optional+Computed with no fixed literal default - "the
	// value returned by Azure depends on when the resource was created if
	// provisioned with `Default`" - so the v4 default isn't confidently
	// known) and no longer accepts "Default" as a value; both changes are
	// registered.
	rules.Register(changedDefaultRule{
		resourceType: "azurerm_mssql_managed_instance",
		attr:         "proxy_override",
		oldDefault:   nil,
		newDefault:   "Redirect",
	})
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_mssql_managed_instance",
		attr:          "proxy_override",
		removedValues: []string{"Default"},
	})

	// azurerm_mssql_managed_instance_transparent_data_encryption: the
	// deprecated managed_hsm_key_id property was removed in favour of
	// key_vault_key_id (consolidation precedent).
	rules.Register(flagAttrRule{
		resourceType: "azurerm_mssql_managed_instance_transparent_data_encryption",
		attr:         "managed_hsm_key_id",
		reason:       "in v5 the key is specified via key_vault_key_id; consolidate managed_hsm_key_id into key_vault_key_id",
	})

	// azurerm_mssql_server: minimum_tls_version no longer accepts Disabled,
	// 1.0 or 1.1 as a value.
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_mssql_server",
		attr:          "minimum_tls_version",
		removedValues: []string{"Disabled", "1.0", "1.1"},
	})

	// azurerm_mssql_server_extended_auditing_policy: same
	// storage_endpoint -> blob_storage_endpoint value-preserving rename as
	// azurerm_mssql_database_extended_auditing_policy above.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_mssql_server_extended_auditing_policy",
		oldAttr:      "storage_endpoint",
		newAttr:      "blob_storage_endpoint",
	})

	// azurerm_mssql_server_security_alert_policy: the deprecated
	// email_account_admins property was removed in favour of
	// email_account_admins_enabled. Unlike azurerm_mssql_database's
	// threat_detection_policy.email_account_admins above, this is verified
	// to be a plain TypeBool on both sides of the rename (both pre- and
	// post-5.0 schema declare it as pluginsdk.TypeBool), so this is a
	// straightforward value-preserving bool rename.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_mssql_server_security_alert_policy",
		oldAttr:      "email_account_admins",
		newAttr:      "email_account_admins_enabled",
	})

	// azurerm_mssql_server_transparent_data_encryption: same managed_hsm_key_id
	// -> key_vault_key_id consolidation as
	// azurerm_mssql_managed_instance_transparent_data_encryption above.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_mssql_server_transparent_data_encryption",
		attr:         "managed_hsm_key_id",
		reason:       "in v5 the key is specified via key_vault_key_id; consolidate managed_hsm_key_id into key_vault_key_id",
	})

	// azurerm_mssql_virtual_machine: the deprecated
	// auto_backup.encryption_enabled property was removed in v5 (encryption
	// is now inferred from whether encryption_password is set).
	rules.Register(removedAttrRule{
		resourceType: "azurerm_mssql_virtual_machine",
		blockPath:    []string{"auto_backup"},
		attr:         "encryption_enabled",
	})

	// misc
	//
	// azurerm_api_management: the deprecated hostname_configuration.*.key_vault_id
	// properties were removed in favour of hostname_configuration.*.key_vault_certificate_id.
	// Verified against source (internal/services/apimanagement/api_management_resource.go):
	// both hold the same Key Vault certificate secret ID value, so these are
	// value-preserving renames.
	for _, hostnameBlock := range []string{"developer_portal", "management", "portal", "proxy", "scm"} {
		rules.Register(renamedAttrRule{
			resourceType: "azurerm_api_management",
			blockPath:    []string{"hostname_configuration", hostnameBlock},
			oldAttr:      "key_vault_id",
			newAttr:      "key_vault_certificate_id",
		})
	}
	// protocols.enable_http2 -> protocols.http2_enabled and the remaining
	// security.enable_* -> security.*_enabled bool renames (verified against
	// source: all are plain value-preserving TypeBool renames, same shape as
	// security.enable_backend_ssl30 -> security.backend_ssl30_enabled already
	// registered above). Note security.enable_backend_ssl30 is intentionally
	// not repeated here to avoid a duplicate rule ID.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_api_management",
		blockPath:    []string{"protocols"},
		oldAttr:      "enable_http2",
		newAttr:      "http2_enabled",
	})
	for _, oldNew := range [][2]string{
		{"enable_backend_tls10", "backend_tls10_enabled"},
		{"enable_backend_tls11", "backend_tls11_enabled"},
		{"enable_frontend_ssl30", "frontend_ssl30_enabled"},
		{"enable_frontend_tls10", "frontend_tls10_enabled"},
		{"enable_frontend_tls11", "frontend_tls11_enabled"},
	} {
		rules.Register(renamedAttrRule{
			resourceType: "azurerm_api_management",
			blockPath:    []string{"security"},
			oldAttr:      oldNew[0],
			newAttr:      oldNew[1],
		})
	}

	// azurerm_api_management_custom_domain: same key_vault_id ->
	// key_vault_certificate_id value-preserving rename as azurerm_api_management
	// above, scoped to each of its host-type blocks.
	for _, hostnameBlock := range []string{"developer_portal", "gateway", "management", "portal", "scm"} {
		rules.Register(renamedAttrRule{
			resourceType: "azurerm_api_management_custom_domain",
			blockPath:    []string{hostnameBlock},
			oldAttr:      "key_vault_id",
			newAttr:      "key_vault_certificate_id",
		})
	}

	// azurerm_application_insights: all three deprecated properties are
	// boolean INVERSIONS, not value-preserving renames. Verified against
	// source (internal/services/applicationinsights/application_insights_resource.go):
	// e.g. `DisableIPMasking = pointer.To(!d.Get("ip_masking_enabled").(bool))`
	// when the new attr is used, vs `DisableIPMasking = pointer.To(d.Get("disable_ip_masking").(bool))`
	// when the old attr is used - disabled=true means the same thing as
	// enabled=false. All three are flagged for manual review rather than
	// auto-renamed.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_application_insights",
		attr:         "daily_data_cap_notifications_disabled",
		reason:       "daily_data_cap_notifications_disabled was removed in v5 in favour of daily_data_cap_notifications_enabled, and the boolean meaning is inverted; update the value manually (e.g. daily_data_cap_notifications_disabled = true becomes daily_data_cap_notifications_enabled = false)",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_application_insights",
		attr:         "disable_ip_masking",
		reason:       "disable_ip_masking was removed in v5 in favour of ip_masking_enabled, and the boolean meaning is inverted; update the value manually (e.g. disable_ip_masking = true becomes ip_masking_enabled = false)",
	})
	rules.Register(flagAttrRule{
		resourceType: "azurerm_application_insights",
		attr:         "local_authentication_disabled",
		reason:       "local_authentication_disabled was removed in v5 in favour of local_authentication_enabled, and the boolean meaning is inverted; update the value manually (e.g. local_authentication_disabled = true becomes local_authentication_enabled = false)",
	})

	// NOTE: azurerm_application_gateway's two v5 changes (enable_http2 ->
	// http2_enabled, ssl_profile.verify_client_cert_issuer_dn ->
	// verify_client_certificate_issuer_dn) are already registered in the
	// networking section above; not duplicated here.

	// azurerm_automation_account: the encryption.key_source property was
	// removed in v5.
	rules.Register(removedAttrRule{
		resourceType: "azurerm_automation_account",
		blockPath:    []string{"encryption"},
		attr:         "key_source",
	})

	// azurerm_batch_pool: the deprecated certificate block was removed in v5
	// with no direct replacement (deprecated by the service itself).
	rules.Register(removedBlockRule{
		resourceType: "azurerm_batch_pool",
		blockPath:    []string{"certificate"},
		reason:       "the deprecated certificate block was removed in v5 as it is no longer supported by the Batch service",
	})

	// azurerm_bot_channel_ms_teams: enable_calling -> calling_enabled.
	// Verified against source (internal/services/bot/bot_channel_ms_teams_resource.go):
	// both are plain TypeBool and map directly to EnableCalling, so this is a
	// value-preserving rename.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_bot_channel_ms_teams",
		oldAttr:      "enable_calling",
		newAttr:      "calling_enabled",
	})

	// azurerm_bot_channels_registration / azurerm_bot_service_azure_bot /
	// azurerm_bot_web_app: microsoft_app_type is now Required in v5 (it was
	// Optional pre-5.0). The tool can't supply a value, so this flags the
	// attribute as needing manual review when it's absent.
	for _, resourceType := range []string{
		"azurerm_bot_channels_registration",
		"azurerm_bot_service_azure_bot",
		"azurerm_bot_web_app",
	} {
		rules.Register(requiredAttrRule{
			resourceType: resourceType,
			attr:         "microsoft_app_type",
			reason:       "microsoft_app_type is required in v5; add it explicitly",
		})
	}

	// azurerm_cdn_endpoint_custom_domain: cdn_managed_https.tls_version and
	// user_managed_https.tls_version no longer accept None or TLS10 as a
	// value; the tool can't safely pick a replacement, so this only flags the
	// value for manual review.
	for _, block := range []string{"cdn_managed_https", "user_managed_https"} {
		rules.Register(valueChangeRule{
			resourceType:  "azurerm_cdn_endpoint_custom_domain",
			blockPath:     []string{block},
			attr:          "tls_version",
			removedValues: []string{"None", "TLS10"},
		})
	}

	// azurerm_cdn_frontdoor_custom_domain: the deprecated
	// tls.minimum_tls_version property was removed in favour of
	// tls.minimum_version. Verified against source
	// (internal/services/cdn/cdn_frontdoor_custom_domain_resource.go): both
	// are plain TypeString holding the same TLS-version value, so this is a
	// value-preserving rename. Separately, tls.minimum_version (under either
	// name) no longer accepts TLS10 as a value in v5, so that's also flagged.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_cdn_frontdoor_custom_domain",
		blockPath:    []string{"tls"},
		oldAttr:      "minimum_tls_version",
		newAttr:      "minimum_version",
	})
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_cdn_frontdoor_custom_domain",
		blockPath:     []string{"tls"},
		attr:          "minimum_version",
		removedValues: []string{"TLS10"},
	})

	// azurerm_communication_service: data_location is now Required in v5 and
	// no longer defaults to "United States" (verified against source:
	// pre-5.0 it's Optional with Default: "United States"; v5 is Required
	// with no default). The tool can't supply a value, so this flags the
	// attribute as needing manual review when it's absent.
	rules.Register(requiredAttrRule{
		resourceType: "azurerm_communication_service",
		attr:         "data_location",
		reason:       "data_location is required in v5 and no longer defaults to \"United States\"; add it explicitly",
	})

	// azurerm_container_app / azurerm_container_app_job: the deprecated
	// template.container.liveness_probe.termination_grace_period_seconds and
	// template.container.startup_probe.termination_grace_period_seconds
	// properties were removed in v5.
	for _, resourceType := range []string{"azurerm_container_app", "azurerm_container_app_job"} {
		for _, probe := range []string{"liveness_probe", "startup_probe"} {
			rules.Register(removedAttrRule{
				resourceType: resourceType,
				blockPath:    []string{"template", "container", probe},
				attr:         "termination_grace_period_seconds",
			})
		}
	}

	// azurerm_container_app_environment: logs_destination is no longer
	// Computed in v5 and now defaults to empty (Streaming Only), instead of
	// its pre-5.0 behavior of being implicitly derived from whether
	// log_analytics_workspace_id was set. Verified against source
	// (internal/services/containerapps/container_app_environment_resource.go):
	// pre-5.0 it's Optional+Computed with no fixed literal default, so the v4
	// default isn't confidently known.
	rules.Register(changedDefaultRule{
		resourceType: "azurerm_container_app_environment",
		attr:         "logs_destination",
		oldDefault:   nil,
		newDefault:   "\"\" (streaming only; must be set to log-analytics to use log_analytics_workspace_id)",
	})

	// azurerm_container_registry: the deprecated trust_policy_enabled
	// property was removed in v5 as it is no longer supported by the
	// service.
	rules.Register(removedAttrRule{
		resourceType: "azurerm_container_registry",
		attr:         "trust_policy_enabled",
	})
	// NOTE: "The encryption block is no longer Computed. It now defaults to
	// empty, meaning encryption will be disabled." is a block-level
	// Computed-ness/default change, not a renamed/removed/required attribute
	// or block - none of the available rule primitives model "a block's
	// implicit default changed" without misrepresenting it as the block
	// being newly required (it isn't; it's optional either way). Skipped as
	// not mechanically encodable with the current primitive catalog.

	// azurerm_cost_management_scheduled_action: validation for email_subject
	// now enforces a maximum length of 50 characters; the tool can't verify
	// string length against an arbitrary/interpolated value, so this is
	// flagged for manual review rather than auto-migrated.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_cost_management_scheduled_action",
		attr:         "email_subject",
		reason:       "email_subject now validates a maximum length of 50 characters in v5; confirm the value fits",
	})

	// azurerm_dashboard_grafana: the deprecated Essential value for sku was
	// removed, and grafana_major_version no longer accepts "11" as a value.
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_dashboard_grafana",
		attr:          "sku",
		removedValues: []string{"Essential"},
	})
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_dashboard_grafana",
		attr:          "grafana_major_version",
		removedValues: []string{"11"},
	})

	// azurerm_data_factory_integration_runtime_self_hosted: validation for
	// rbac_authorization.resource_id now requires an integration runtime
	// resource ID (case-sensitive) rather than a non-empty string; the tool
	// can't verify shape/casing on the caller's behalf, so this is flagged
	// for manual review rather than auto-migrated.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_data_factory_integration_runtime_self_hosted",
		blockPath:    []string{"rbac_authorization"},
		attr:         "resource_id",
		reason:       "rbac_authorization.resource_id now validates its value as an integration runtime resource ID (case-sensitive) in v5; confirm the value is a correctly-cased integration runtime resource ID",
	})

	// azurerm_data_factory_linked_service_azure_blob_storage: the deprecated
	// key_vault_sas_token block was removed in favour of the
	// sas_token_linked_key_vault_key block. Verified against source
	// (internal/services/datafactory/data_factory_linked_service_azure_blob_storage_resource.go):
	// both blocks declare the identical nested schema (linked_service_name,
	// secret_name), but a block-to-block rename can't be expressed with the
	// attribute-only renamedAttrRule primitive, so this is flagged for
	// manual review instead.
	rules.Register(removedBlockRule{
		resourceType: "azurerm_data_factory_linked_service_azure_blob_storage",
		blockPath:    []string{"key_vault_sas_token"},
		reason:       "key_vault_sas_token was renamed to sas_token_linked_key_vault_key in v5 (identical nested schema: linked_service_name/secret_name); rename the block manually",
	})

	// azurerm_data_factory_linked_service_azure_databricks: the deprecated
	// msi_work_space_resource_id property was removed in favour of
	// msi_workspace_id. Verified against source: both hold the same
	// Workspace Resource ID value, so this is a value-preserving rename.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_data_factory_linked_service_azure_databricks",
		oldAttr:      "msi_work_space_resource_id",
		newAttr:      "msi_workspace_id",
	})

	// azurerm_data_factory_linked_service_mysql: driver_version now defaults
	// to V2 (verified against source: pre-5.0 default is V1).
	dataFactoryMySQLDriverVersionOldDefault := cty.StringVal("V1")
	rules.Register(changedDefaultRule{
		resourceType: "azurerm_data_factory_linked_service_mysql",
		attr:         "driver_version",
		oldDefault:   &dataFactoryMySQLDriverVersionOldDefault,
		newDefault:   "V2",
	})

	// azurerm_datadog_monitor_sso_configuration: the deprecated
	// single_sign_on_enabled property was removed in favour of
	// single_sign_on. Verified against source
	// (internal/services/datadog/datadog_monitor_sso_configuration_resource.go):
	// despite the "_enabled" suffix, single_sign_on_enabled is a plain
	// TypeString enum (Enable/Disable) - the same type as single_sign_on -
	// both assign directly to SingleSignOnState with no inversion, so this
	// is a value-preserving rename, not a boolean inversion.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_datadog_monitor_sso_configuration",
		oldAttr:      "single_sign_on_enabled",
		newAttr:      "single_sign_on",
	})

	// azurerm_data_factory_pipeline: the deprecated moniter_metrics_after_duration
	// property (a typo in the original attribute name) was removed in favour
	// of the correctly-spelled monitor_metrics_after_duration property; same
	// value, value-preserving rename.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_data_factory_pipeline",
		oldAttr:      "moniter_metrics_after_duration",
		newAttr:      "monitor_metrics_after_duration",
	})

	// azurerm_iot_security_solution: the deprecated recommendations_enabled
	// block was removed in favour of the recommendations block. Verified
	// against source (v4.50.0 tag vs main,
	// internal/services/securitycenter/iot_security_solution_resource.go):
	// both declare an identical nested schema (acr_authentication,
	// agent_send_unutilized_msg, baseline, etc.), but a block-to-block
	// rename can't be expressed with the attribute-only renamedAttrRule
	// primitive, so this is flagged for manual review instead.
	rules.Register(removedBlockRule{
		resourceType: "azurerm_iot_security_solution",
		blockPath:    []string{"recommendations_enabled"},
		reason:       "recommendations_enabled was renamed to recommendations in v5 (identical nested schema); rename the block manually",
	})

	// azurerm_kusto_attached_database_configuration: the deprecated
	// cluster_resource_id property was removed in favour of cluster_id. Both
	// hold the same cluster resource ID value, so this is a value-preserving
	// rename.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_kusto_attached_database_configuration",
		oldAttr:      "cluster_resource_id",
		newAttr:      "cluster_id",
	})

	// azurerm_kusto_cluster: the deprecated language_extensions block was
	// removed in favour of language_extension. Verified against source
	// (internal/services/kusto/kusto_cluster_resource.go): both declare the
	// identical nested schema (name, image), but a block-to-block rename
	// can't be expressed with the attribute-only renamedAttrRule primitive,
	// so this is flagged for manual review instead. Separately, the
	// deprecated virtual_network_configuration block was removed entirely in
	// v5 (no longer supported by Azure), with no replacement.
	rules.Register(removedBlockRule{
		resourceType: "azurerm_kusto_cluster",
		blockPath:    []string{"language_extensions"},
		reason:       "language_extensions was renamed to language_extension in v5 (identical nested schema: name/image); rename the block manually",
	})
	rules.Register(removedBlockRule{
		resourceType: "azurerm_kusto_cluster",
		blockPath:    []string{"virtual_network_configuration"},
		reason:       "virtual_network_configuration was removed in v5 as it is no longer supported by Azure, with no replacement",
	})

	// azurerm_kusto_eventgrid_data_connection: the deprecated
	// eventgrid_resource_id and managed_identity_resource_id properties were
	// removed in favour of eventgrid_event_subscription_id and
	// managed_identity_id respectively. Both pairs hold the same shape of ID
	// value, so these are value-preserving renames.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_kusto_eventgrid_data_connection",
		oldAttr:      "eventgrid_resource_id",
		newAttr:      "eventgrid_event_subscription_id",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_kusto_eventgrid_data_connection",
		oldAttr:      "managed_identity_resource_id",
		newAttr:      "managed_identity_id",
	})

	// NOTE: azurerm_lb's only documented v5 change ("The subnet_id property
	// is no longer Computed. The public_ip_address_id property is no longer
	// Computed.") is a computed-behavior change only (no config-visible
	// attribute rename, removal, or value change); skipped, same reasoning
	// as azurerm_mssql_database.enclave_type above.

	// azurerm_local_network_gateway: address_space's type changed from a
	// List to a Set in v5, making it unordered; index-based references into
	// it are no longer possible. This is a type change, not a mechanical
	// rewrite, so it's flagged for manual review.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_local_network_gateway",
		attr:         "address_space",
		reason:       "address_space's type changed from a List to a Set in v5 (now unordered); review any index-based references to its elements",
	})

	// azurerm_log_analytics_linked_storage_account: the deprecated
	// workspace_resource_id property was removed in favour of workspace_id.
	// Both hold the same Log Analytics Workspace ID value, so this is a
	// value-preserving rename.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_log_analytics_linked_storage_account",
		oldAttr:      "workspace_resource_id",
		newAttr:      "workspace_id",
	})

	// NOTE: azurerm_log_analytics_storage_insights's only documented v5
	// change ("The storage_account_id property is marked as ForceNew.") is a
	// plan-time behavior change only (no config-visible attribute rename,
	// removal, or value change); skipped, same reasoning as azurerm_lb above.

	// azurerm_log_analytics_workspace: the deprecated
	// local_authentication_disabled property was removed in v5 in favour of
	// local_authentication_enabled. Same INVERSION pattern as
	// azurerm_cosmosdb_account.local_authentication_disabled and
	// azurerm_application_insights.local_authentication_disabled above (not
	// value-preserving), so this is flagged for manual review rather than
	// auto-renamed.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_log_analytics_workspace",
		attr:         "local_authentication_disabled",
		reason:       "local_authentication_disabled was removed in v5 in favour of local_authentication_enabled, and the boolean meaning is inverted; update the value manually (e.g. local_authentication_disabled = true becomes local_authentication_enabled = false)",
	})

	// azurerm_maintenance_assignment_virtual_machine_scale_set:
	// virtual_machine_scale_set_id is no longer case-insensitive in v5; the
	// tool can't verify casing against the caller's actual resource ID, so
	// this is flagged for manual review rather than auto-migrated.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_maintenance_assignment_virtual_machine_scale_set",
		attr:         "virtual_machine_scale_set_id",
		reason:       "virtual_machine_scale_set_id is no longer case-insensitive in v5; confirm the value's casing matches the actual resource ID exactly",
	})

	// azurerm_monitor_aad_diagnostic_setting / azurerm_monitor_diagnostic_setting:
	// the deprecated enabled_log.retention_policy block was removed in v5 in
	// favour of the azurerm_storage_management_policy resource.
	for _, resourceType := range []string{"azurerm_monitor_aad_diagnostic_setting", "azurerm_monitor_diagnostic_setting"} {
		rules.Register(removedBlockRule{
			resourceType: resourceType,
			blockPath:    []string{"enabled_log", "retention_policy"},
			reason:       "the deprecated enabled_log.retention_policy block was removed in v5; manage log retention via the azurerm_storage_management_policy resource instead",
		})
	}
	// azurerm_monitor_diagnostic_setting only: the deprecated metric block
	// was removed in favour of the enabled_metric block, and its nested
	// metric.retention_policy block was likewise removed. Both are
	// registered per the guide, even though flagging metric as removed
	// already implies its retention_policy sub-block is gone too.
	rules.Register(removedBlockRule{
		resourceType: "azurerm_monitor_diagnostic_setting",
		blockPath:    []string{"metric"},
		reason:       "the deprecated metric block was removed in v5 in favour of the enabled_metric block; migrate manually",
	})
	rules.Register(removedBlockRule{
		resourceType: "azurerm_monitor_diagnostic_setting",
		blockPath:    []string{"metric", "retention_policy"},
		reason:       "the deprecated metric.retention_policy block was removed in v5; manage log retention via the azurerm_storage_management_policy resource instead",
	})

	// azurerm_mysql_flexible_server: the deprecated
	// customer_managed_key.managed_hsm_key_id property was removed in favour
	// of customer_managed_key.key_vault_key_id (consolidation precedent, same
	// as azurerm_storage_account.customer_managed_key.managed_hsm_key_id
	// above).
	rules.Register(flagAttrRule{
		resourceType: "azurerm_mysql_flexible_server",
		blockPath:    []string{"customer_managed_key"},
		attr:         "managed_hsm_key_id",
		reason:       "in v5 the key is specified via key_vault_key_id; consolidate managed_hsm_key_id into key_vault_key_id",
	})
	// NOTE: the deprecated public_network_access_enabled property is
	// verified against source (internal/services/mysql/mysql_flexible_server_resource.go)
	// to be Computed-only pre-5.0 (Type: TypeBool, Computed: true, no
	// Optional) - i.e. it was already read-only, not user-settable - so its
	// removal in v5 in favour of public_network_access is a computed-only
	// change with nothing for the tool to edit or flag; skipped.

	// azurerm_netapp_volume: the deprecated export_policy_rule.protocols_enabled
	// property was removed in favour of export_policy_rule.protocol.
	// Verified against source (internal/services/netapp/netapp_volume_resource.go):
	// both are an identically-shaped TypeList of protocol name strings
	// (MinItems 1, MaxItems 1, same allowed values), so this is a
	// value-preserving rename.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_netapp_volume",
		blockPath:    []string{"export_policy_rule"},
		oldAttr:      "protocols_enabled",
		newAttr:      "protocol",
	})

	// azurerm_nginx_deployment: the deprecated diagnose_support_enabled and
	// managed_resource_group properties were removed in v5, and the
	// deprecated logging_storage_account block was removed in favour of the
	// azurerm_monitor_diagnostic_setting resource.
	rules.Register(removedAttrRule{
		resourceType: "azurerm_nginx_deployment",
		attr:         "diagnose_support_enabled",
	})
	rules.Register(removedBlockRule{
		resourceType: "azurerm_nginx_deployment",
		blockPath:    []string{"logging_storage_account"},
		reason:       "the deprecated logging_storage_account block was removed in v5; manage diagnostics via the azurerm_monitor_diagnostic_setting resource instead",
	})
	rules.Register(removedAttrRule{
		resourceType: "azurerm_nginx_deployment",
		attr:         "managed_resource_group",
	})

	// azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack /
	// _virtual_hub_panorama / _virtual_network_local_rulestack /
	// _virtual_network_panorama: plan_id now defaults to panw-cngfw-payg in
	// v5 (verified against source, internal/services/paloalto/*.go: all four
	// resources share the identical pattern of a pre-5.0 default of
	// panw-cloud-ngfw-payg overriding the v5 base default of
	// panw-cngfw-payg).
	paloAltoPlanIDOldDefault := cty.StringVal("panw-cloud-ngfw-payg")
	for _, resourceType := range []string{
		"azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack",
		"azurerm_palo_alto_next_generation_firewall_virtual_hub_panorama",
		"azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack",
		"azurerm_palo_alto_next_generation_firewall_virtual_network_panorama",
	} {
		rules.Register(changedDefaultRule{
			resourceType: resourceType,
			attr:         "plan_id",
			oldDefault:   &paloAltoPlanIDOldDefault,
			newDefault:   "panw-cngfw-payg",
		})
	}

	// azurerm_policy_set_definition: the management_group_id property was
	// removed in v5 in favour of the azurerm_management_group_policy_set_definition
	// resource.
	rules.Register(removedAttrRule{
		resourceType: "azurerm_policy_set_definition",
		attr:         "management_group_id",
	})

	// azurerm_powerbi_embedded: mode now defaults to Gen2 (verified against
	// source, internal/services/powerbi/powerbi_embedded_resource.go:
	// pre-5.0 default is Gen1).
	powerBIEmbeddedModeOldDefault := cty.StringVal("Gen1")
	rules.Register(changedDefaultRule{
		resourceType: "azurerm_powerbi_embedded",
		attr:         "mode",
		oldDefault:   &powerBIEmbeddedModeOldDefault,
		newDefault:   "Gen2",
	})

	// azurerm_redis_cache: minimum_tls_version no longer accepts 1.0 or 1.1
	// as a value.
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_redis_cache",
		attr:          "minimum_tls_version",
		removedValues: []string{"1.0", "1.1"},
	})

	// azurerm_sentinel_alert_rule_fusion: the deprecated name property was
	// removed in v5.
	rules.Register(removedAttrRule{
		resourceType: "azurerm_sentinel_alert_rule_fusion",
		attr:         "name",
	})

	// NOTE: azurerm_sentinel_threat_intelligence_indicator's only documented v5
	// change ("Changing the description, display_name, or validate_from_utc
	// properties now forces a new resource to be created.") is a plan-time
	// ForceNew behavior change only (no config-visible attribute rename,
	// removal, or value change); skipped, same reasoning as azurerm_lb and
	// azurerm_log_analytics_storage_insights above.

	// azurerm_site_recovery_protection_container_mapping:
	// automatic_update.automation_account_id is now Required in v5 (verified
	// against source: pre-5.0 it's Optional). The tool can't supply a value,
	// so this flags the attribute as needing manual review when it's absent.
	// Separately, automatic_update.enabled (a pre-5.0 TypeBool, Optional,
	// Default false) was removed entirely in v5 - to disable automatic
	// update, omit the automatic_update block instead.
	rules.Register(requiredAttrRule{
		resourceType: "azurerm_site_recovery_protection_container_mapping",
		blockPath:    []string{"automatic_update"},
		attr:         "automation_account_id",
		reason:       "automatic_update.automation_account_id is required in v5; add it explicitly",
	})
	rules.Register(removedAttrRule{
		resourceType: "azurerm_site_recovery_protection_container_mapping",
		blockPath:    []string{"automatic_update"},
		attr:         "enabled",
	})

	// azurerm_spring_cloud_connection: client_type no longer accepts "none"
	// as a value (omitting the property sets "none" in the request).
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_spring_cloud_connection",
		attr:          "client_type",
		removedValues: []string{"none"},
	})

	// azurerm_static_web_app: identity.type no longer accepts
	// "SystemAssigned, UserAssigned" as a value.
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_static_web_app",
		blockPath:     []string{"identity"},
		attr:          "type",
		removedValues: []string{"SystemAssigned, UserAssigned"},
	})

	// NOTE: azurerm_storage_container / azurerm_storage_queue /
	// azurerm_storage_share's deprecated resource_manager_id property (removed
	// in favour of id) is verified against source
	// (internal/services/storage/storage_{container,queue,share}_resource.go)
	// to be Computed-only pre-5.0 (Type: TypeString, Computed: true) on all
	// three resources - i.e. it was already read-only, not user-settable -
	// so its removal in v5 is a computed-only change with nothing for the
	// tool to edit or flag; skipped.

	// azurerm_synapse_spark_pool: spark_version no longer accepts 3.2 or 3.3
	// as a value.
	rules.Register(valueChangeRule{
		resourceType:  "azurerm_synapse_spark_pool",
		attr:          "spark_version",
		removedValues: []string{"3.2", "3.3"},
	})

	// === v5 coverage sub-project #2: resource-attribute gaps ===

	// azurerm_log_analytics_workspace: internet_ingestion_enabled and
	// internet_query_enabled become enum-string *_access_type fields
	// (true -> "Enabled", false -> "Disabled").
	rules.Register(boolToEnumRule{
		resourceType: "azurerm_log_analytics_workspace",
		oldAttr:      "internet_ingestion_enabled",
		newAttr:      "internet_ingestion_access_type",
		trueValue:    "Enabled",
		falseValue:   "Disabled",
	})
	rules.Register(boolToEnumRule{
		resourceType: "azurerm_log_analytics_workspace",
		oldAttr:      "internet_query_enabled",
		newAttr:      "internet_query_access_type",
		trueValue:    "Enabled",
		falseValue:   "Disabled",
	})

	// azurerm_application_gateway: the authentication_certificate block (top-level)
	// and backend_http_settings.authentication_certificate block are removed with
	// no mechanical replacement (manual migration to trusted_root_certificate).
	rules.Register(removedBlockRule{
		resourceType: "azurerm_application_gateway",
		blockPath:    []string{"authentication_certificate"},
		reason:       "removed in v5; migrate to trusted_root_certificate / trusted_root_certificate_names",
	})
	rules.Register(removedBlockRule{
		resourceType: "azurerm_application_gateway",
		blockPath:    []string{"backend_http_settings", "authentication_certificate"},
		reason:       "removed in v5; migrate to trusted_root_certificate_names on the backend_http_settings",
	})

	// azurerm_container_registry: georeplications.regional_endpoint_enabled
	// renamed to global_endpoint_routing_enabled (same polarity), now Required.
	// NOTE: register the rename before the now-required rule below - the engine
	// runs rules in registration order over one shared tree, so renaming
	// regional_endpoint_enabled first means the now-required check sees
	// global_endpoint_routing_enabled present and stays silent.
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_container_registry",
		blockPath:    []string{"georeplications"},
		oldAttr:      "regional_endpoint_enabled",
		newAttr:      "global_endpoint_routing_enabled",
	})
	rules.Register(requiredAttrRule{
		resourceType: "azurerm_container_registry",
		blockPath:    []string{"georeplications"},
		attr:         "global_endpoint_routing_enabled",
		reason:       "each georeplications block must set it explicitly",
	})

	// azurerm_hdinsight_*_cluster: storage_account/gen2 ID renames, the
	// storage_container_id -> storage_container_url rename that also requires a
	// Data-Plane URL value (flag), and tls_min_version now required.
	for _, hdi := range []string{
		"azurerm_hdinsight_hadoop_cluster",
		"azurerm_hdinsight_hbase_cluster",
		"azurerm_hdinsight_kafka_cluster",
		"azurerm_hdinsight_spark_cluster",
		"azurerm_hdinsight_interactive_query_cluster",
	} {
		rules.Register(renamedAttrRule{
			resourceType: hdi,
			blockPath:    []string{"storage_account"},
			oldAttr:      "storage_resource_id",
			newAttr:      "storage_account_id",
		})
		rules.Register(renamedAttrRule{
			resourceType: hdi,
			blockPath:    []string{"storage_account_gen2"},
			oldAttr:      "storage_resource_id",
			newAttr:      "storage_account_id",
		})
		rules.Register(renamedAttrRule{
			resourceType: hdi,
			blockPath:    []string{"storage_account_gen2"},
			oldAttr:      "managed_identity_resource_id",
			newAttr:      "user_assigned_identity_id",
		})
		rules.Register(flagAttrRule{
			resourceType: hdi,
			blockPath:    []string{"storage_account"},
			attr:         "storage_container_id",
			reason:       "renamed to storage_container_url and the value must be the container's Data-Plane URL (not an ARM resource ID); v5 validates this strictly, so rename the key and confirm the value",
		})
		rules.Register(requiredAttrRule{
			resourceType: hdi,
			attr:         "tls_min_version",
			reason:       "v4 supplied no default, so set it explicitly",
		})
	}

	// azurerm_netapp_volume: mount_ip_addresses (a computed output) is replaced by
	// mount_target (a nested object list). Nothing to rewrite in the block; flag if
	// present so references (mount_ip_addresses[*] -> mount_target[*].ip_address)
	// get a human's attention. Reference rewriting is out of scope.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_netapp_volume",
		attr:         "mount_ip_addresses",
		reason:       "removed in v5 in favour of mount_target (a list of objects); update references to mount_target[*].ip_address",
	})

	// azurerm_cdn_frontdoor_security_policy: cdn_frontdoor_domain_id is now
	// validated case-sensitively; confirm the resource ID segments use canonical
	// casing. No mechanical rewrite is possible.
	rules.Register(flagAttrRule{
		resourceType: "azurerm_cdn_frontdoor_security_policy",
		blockPath:    []string{"security_policies", "firewall", "association", "domain"},
		attr:         "cdn_frontdoor_domain_id",
		reason:       "now validated case-sensitively in v5; confirm the resource ID's segment casing (resourceGroups, providers, etc.) is canonical",
	})

	// === v5 coverage sub-project #3a: restructures ===

	// azurerm_subnet: service_endpoints (list) -> repeated service_endpoint blocks.
	rules.Register(attrToBlockRule{
		resourceType: "azurerm_subnet",
		attr:         "service_endpoints",
		newBlock:     "service_endpoint",
		valueAttr:    "service",
	})
	// azurerm_virtual_network: subnet.service_endpoints -> subnet.service_endpoint blocks.
	rules.Register(attrToBlockRule{
		resourceType: "azurerm_virtual_network",
		blockPath:    []string{"subnet"},
		attr:         "service_endpoints",
		newBlock:     "service_endpoint",
		valueAttr:    "service",
	})

	// azurerm_site_recovery_replicated_vm: network_interface.* IP fields move
	// into a network_interface.ip_configuration block. name is required only
	// when a network_interface has multiple ip_configuration blocks.
	rules.Register(fieldMoveToNestedBlockRule{
		resourceType: "azurerm_site_recovery_replicated_vm",
		parentPath:   []string{"network_interface"},
		childBlock:   "ip_configuration",
		attrs: []string{
			"failover_test_static_ip",
			"failover_test_subnet_name",
			"failover_test_public_ip_address_id",
			"target_static_ip",
			"target_subnet_name",
			"recovery_load_balancer_backend_address_pool_ids",
			"recovery_public_ip_address_id",
		},
		requiredChildAttrsWhenMultiple: []requiredChildAttr{
			{name: "name", reason: "each ip_configuration must be named (and exactly one marked primary)"},
		},
	})

	// === v5 coverage sub-project #3b: cdn_frontdoor_rule ===

	// top-level attribute rename
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		oldAttr:      "behavior_on_match",
		newAttr:      "behaviour_on_match",
	})

	// action block renames (under actions {})
	rules.Register(renamedBlockRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		parentPath:   []string{"actions"},
		oldBlock:     "request_header_action",
		newBlock:     "modify_request_header",
		innerRenames: map[string]string{"header_action": "operator", "value": "header_value"},
	})
	rules.Register(renamedBlockRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		parentPath:   []string{"actions"},
		oldBlock:     "response_header_action",
		newBlock:     "modify_response_header",
		innerRenames: map[string]string{"header_action": "operator", "value": "header_value"},
	})
	rules.Register(renamedBlockRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		parentPath:   []string{"actions"},
		oldBlock:     "url_redirect_action",
		newBlock:     "url_redirect",
		innerRenames: map[string]string{"destination_hostname": "destination_host_name"},
	})
	rules.Register(renamedBlockRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		parentPath:   []string{"actions"},
		oldBlock:     "url_rewrite_action",
		newBlock:     "url_rewrite",
		innerRenames: map[string]string{"destination": "destination_path", "preserve_unmatched_path": "preserve_unmatched_path_enabled"},
	})

	// route_configuration_override_action fan-out. Order matters: rename the
	// block, then rename the flat cache_* fields (BEFORE they move - moved
	// fields become unstructured tokens and can't be renamed afterward), then
	// move fields into the caching / origin_group sub-blocks.
	// NOTE: the now-required caching.behaviour / origin_group.* fields are NOT
	// flagged: moved fields are GetAttribute-invisible, so a requiredAttrRule
	// would false-fire when they ARE supplied. Documented intentional gap.
	//
	// This is the first renamedBlockRule call site whose renamed block is
	// re-addressed by NAME in later rules within the same pass (the other 4
	// registrations only rename attributes inside the block during their own
	// Apply). That surfaced a bug in hashicorp/hcl/v2 v2.24.0's
	// (*hclwrite.Block).SetType - see the RenameBlockType doc comment in
	// internal/rules/helpers.go - which renamedBlockRule now uses instead of
	// SetType so that a block's new type name is visible to later
	// WalkNestedBlocks lookups, not just to the final serialized file.
	rules.Register(renamedBlockRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		parentPath:   []string{"actions"},
		oldBlock:     "route_configuration_override_action",
		newBlock:     "route_configuration_override",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		blockPath:    []string{"actions", "route_configuration_override"},
		oldAttr:      "cache_behavior",
		newAttr:      "behaviour",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		blockPath:    []string{"actions", "route_configuration_override"},
		oldAttr:      "cache_duration",
		newAttr:      "duration",
	})
	rules.Register(renamedAttrRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		blockPath:    []string{"actions", "route_configuration_override"},
		oldAttr:      "query_string_caching_behavior",
		newAttr:      "query_string_behaviour",
	})
	rules.Register(fieldMoveToNestedBlockRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		parentPath:   []string{"actions", "route_configuration_override"},
		childBlock:   "caching",
		attrs:        []string{"behaviour", "duration", "compression_enabled", "query_string_behaviour", "query_string_parameters"},
	})
	rules.Register(fieldMoveToNestedBlockRule{
		resourceType: "azurerm_cdn_frontdoor_rule",
		parentPath:   []string{"actions", "route_configuration_override"},
		childBlock:   "origin_group",
		attrs:        []string{"cdn_frontdoor_origin_group_id", "forwarding_protocol"},
	})

	// cdn_frontdoor_rule conditions: 19 *_condition blocks renamed (each with
	// match_values -> values, some with an extra *_name -> name). negate_condition
	// is folded into operator (new_operator = (negate ? "Not" : "") + operator)
	// when both are literals, or flagged NeedsReview when they aren't.
	// operator becomes required on several, and no longer accepts "Any" on
	// remote_address/socket_address. values (was match_values) becomes required
	// on device_type and remote_address.
	type cdnCond struct {
		oldBlock      string
		newBlock      string
		extraRename   map[string]string // beyond match_values -> values
		operatorReq   bool
		operatorNoAny bool
		valuesReq     bool // values (was match_values) is now required on this condition
	}
	for _, c := range []cdnCond{
		{oldBlock: "client_port_condition", newBlock: "client_port"},
		{oldBlock: "cookies_condition", newBlock: "request_cookies", extraRename: map[string]string{"cookie_name": "name"}},
		{oldBlock: "host_name_condition", newBlock: "host_name"},
		{oldBlock: "http_version_condition", newBlock: "http_version", operatorReq: true},
		{oldBlock: "is_device_condition", newBlock: "device_type", operatorReq: true, valuesReq: true},
		{oldBlock: "post_args_condition", newBlock: "post_argument", extraRename: map[string]string{"post_args_name": "name"}, operatorReq: true},
		{oldBlock: "query_string_condition", newBlock: "query_string"},
		{oldBlock: "remote_address_condition", newBlock: "remote_address", operatorReq: true, operatorNoAny: true, valuesReq: true},
		{oldBlock: "request_body_condition", newBlock: "request_body"},
		{oldBlock: "request_header_condition", newBlock: "request_header", extraRename: map[string]string{"header_name": "name"}},
		{oldBlock: "request_method_condition", newBlock: "request_method", operatorReq: true},
		{oldBlock: "request_scheme_condition", newBlock: "request_scheme", operatorReq: true},
		{oldBlock: "request_uri_condition", newBlock: "request_url"},
		{oldBlock: "server_port_condition", newBlock: "server_port"},
		{oldBlock: "socket_address_condition", newBlock: "socket_address", operatorReq: true, operatorNoAny: true},
		{oldBlock: "ssl_protocol_condition", newBlock: "ssl_protocol", operatorReq: true},
		{oldBlock: "url_file_extension_condition", newBlock: "request_file_extension"},
		{oldBlock: "url_filename_condition", newBlock: "request_filename"},
		{oldBlock: "url_path_condition", newBlock: "request_path"},
	} {
		inner := map[string]string{"match_values": "values"}
		for k, v := range c.extraRename {
			inner[k] = v
		}
		rules.Register(renamedBlockRule{
			resourceType: "azurerm_cdn_frontdoor_rule",
			parentPath:   []string{"conditions"},
			oldBlock:     c.oldBlock,
			newBlock:     c.newBlock,
			innerRenames: inner,
		})
		var skipOperatorValues []string
		if c.operatorNoAny {
			skipOperatorValues = []string{"Any"}
		}
		rules.Register(negateFoldRule{
			resourceType:       "azurerm_cdn_frontdoor_rule",
			blockPath:          []string{"conditions", c.newBlock},
			skipOperatorValues: skipOperatorValues,
		})
		if c.operatorReq {
			rules.Register(requiredAttrRule{
				resourceType: "azurerm_cdn_frontdoor_rule",
				blockPath:    []string{"conditions", c.newBlock},
				attr:         "operator",
				reason:       "the v5 CDN FrontDoor rule schema requires an explicit operator on this condition",
			})
		}
		if c.valuesReq {
			rules.Register(requiredAttrRule{
				resourceType: "azurerm_cdn_frontdoor_rule",
				blockPath:    []string{"conditions", c.newBlock},
				attr:         "values",
				reason:       "the v5 CDN FrontDoor rule schema requires values on this condition",
			})
		}
		if c.operatorNoAny {
			rules.Register(valueChangeRule{
				resourceType:  "azurerm_cdn_frontdoor_rule",
				blockPath:     []string{"conditions", c.newBlock},
				attr:          "operator",
				removedValues: []string{"Any"},
			})
		}
	}

	// === v5 coverage sub-project #4: data sources ===
	// NOTE: most guide data-source changes are computed-OUTPUT-only (enable_X->X_enabled
	// renames, *_arm_resource_id renames, resource_manager_id->id, etc.) - nothing to
	// rewrite in the data block, references are out of scope - so they are intentionally
	// skipped. Only INPUT-argument changes are handled below.

	// Storage data sources: storage_account_name (input) -> storage_account_id.
	for _, ds := range []string{
		"azurerm_storage_container",
		"azurerm_storage_queue",
		"azurerm_storage_share",
		"azurerm_storage_table",
	} {
		rules.Register(nameToIDRule{
			blockKind:    rules.DataBlock,
			resourceType: ds,
			oldAttr:      "storage_account_name",
			newAttr:      "storage_account_id",
			targetType:   "azurerm_storage_account",
		})
	}

	// Data-source input args that collapse into a single parent-resource ID that
	// can't be mechanically constructed - flag for manual migration.
	rules.Register(flagAttrRule{blockKind: rules.DataBlock, resourceType: "azurerm_storage_blob", attr: "storage_account_name", reason: "removed in v5; the storage_blob data source now uses storage_container_id"})
	rules.Register(flagAttrRule{blockKind: rules.DataBlock, resourceType: "azurerm_storage_blob", attr: "storage_container_name", reason: "removed in v5; the storage_blob data source now uses storage_container_id"})

	for _, sb := range []struct{ ds, reason string }{
		{"azurerm_servicebus_namespace_disaster_recovery_config", "removed in v5; use namespace_id"},
		{"azurerm_servicebus_queue", "removed in v5; use namespace_id"},
	} {
		rules.Register(flagAttrRule{blockKind: rules.DataBlock, resourceType: sb.ds, attr: "namespace_name", reason: sb.reason})
		rules.Register(flagAttrRule{blockKind: rules.DataBlock, resourceType: sb.ds, attr: "resource_group_name", reason: sb.reason})
	}
	rules.Register(flagAttrRule{blockKind: rules.DataBlock, resourceType: "azurerm_servicebus_subscription", attr: "namespace_name", reason: "removed in v5; use topic_id"})
	rules.Register(flagAttrRule{blockKind: rules.DataBlock, resourceType: "azurerm_servicebus_subscription", attr: "resource_group_name", reason: "removed in v5; use topic_id"})
	rules.Register(flagAttrRule{blockKind: rules.DataBlock, resourceType: "azurerm_servicebus_subscription", attr: "topic_name", reason: "removed in v5; use topic_id"})

	rules.Register(flagAttrRule{blockKind: rules.DataBlock, resourceType: "azurerm_storage_table_entity", attr: "storage_table_id", reason: "v5 no longer accepts the legacy Data-Plane URL form; supply the management-plane storage table resource ID"})
}
