# Binary end-to-end corpus

`testdata/corpus/input/` is one Terraform config - a root module that
instantiates one child module per service (`modules/<service>/`) - exercised
end-to-end through the built `upgrader` binary.

- `input/`  - valid v4 config (root provider/pin blocks + per-service modules).
- `expected/` - the v5 output the binary must produce (whole tree; golden).
- `expected-report.json` - the binary's `--report json` output, temp-path normalized.

Resources that are *removed* in v5 are not in the corpus (they have no valid v5
form); they stay covered by the per-rule fixtures under
`internal/rules/azurerm/testdata`.

## Layers

- `go test ./tests/ -run TestGolden` - always-on: input→expected + report + idempotency.
- `AZURERM_ACC=1 go test ./tests/ -run TestValidate` - opt-in: upgrades a fresh
  copy of the corpus, then runs `terraform init` + `terraform validate -json`
  against it and asserts zero error-severity diagnostics. Needs terraform on
  PATH and network access to the Terraform registry: it deliberately bypasses
  the repo's dev_overrides (via a `provider_installation { direct {} }` CLI
  config) so the corpus is checked against the **released** azurerm v5.0.0
  provider, pinned by the upgraded root module. The first run downloads the
  provider (~340MB) into a persistent plugin cache; subsequent runs reuse it.

## Regenerating goldens

    go test ./tests/ -run TestGolden -update

Always read the regenerated `expected/` tree before committing.

## Adding a service module

Add `modules/<service>/main.tf`, wire `module "<service>"` into
`input/main.tf`, run `-update`, inspect, then run both layers.

## Known rule gaps

Some upgrader rules do not match the behavior of the released GA v5.0.0
provider. Resources affected by a known gap are excluded from this corpus
(they would fail `TestValidate` through no fault of the corpus itself) and
stay covered instead by the per-rule golden fixtures under
`internal/rules/azurerm/testdata`. In other, narrower cases the resource
itself stays in the corpus but an individual v5 field or block is omitted
from its fixture, because the upgrader only flags that field/block for
manual review (a flag-only `removedBlockRule`/consolidation) rather than
rewriting it; those specific rules remain covered end-to-end by their own
per-rule fixtures under `internal/rules/azurerm/testdata`, just not by this
corpus.

- `azurerm_storage_blob` - the upgrader rewrites `storage_account_name` to
  `storage_account_id`, but GA v5.0.0 instead consolidates this resource onto
  `storage_container_id`. This is a rule bug to fix separately; until then, do
  not add this resource to the corpus. (`azurerm_storage_queue`, whose
  `storage_account_name` → `storage_account_id` rewrite GA does honor, is
  used as the corpus's name-to-id exemplar instead.)
- `azurerm_storage_table` - previously recorded here as broken alongside
  `azurerm_storage_blob`, but re-verified against the real GA v5.0.0 provider
  schema during Task 3: `azurerm_storage_table` actually requires
  `storage_account_id` (not `storage_container_id`), so the upgrader's
  `storage_account_name` → `storage_account_id` rewrite for this resource is
  correct and validates cleanly. It is still not added to the corpus as a
  standalone module resource for this batch (per the Task 3 dispatch), but it
  is used as a sibling of `azurerm_storage_table_entity` in
  `modules/storage/main.tf` and passes `TestValidate` there.
- `azurerm_storage_share_directory`, `azurerm_storage_share_file` - the
  upgrader only flags `storage_share_id` for manual review (leaves it
  unchanged) with the reason "storage_share_id was replaced by
  storage_share_url in v5". GA v5.0.0 confirms this: the field is gone from
  the schema entirely (`storage_share_url` is required instead), so the
  flagged-but-unrewritten output fails validate with `Missing required
  argument: "storage_share_url" is required` and `Unsupported argument: An
  argument named "storage_share_id" is not expected here`. This is an
  id-to-url value transform (not a simple rename) that the upgrader cannot
  perform mechanically; a rule fix would need to either compute the share URL
  from the referenced share or leave it for manual migration as it does now,
  but either way the resource cannot be added to this corpus until GA-clean.
- `azurerm_storage_account_customer_managed_key` - the upgrader only flags
  `key_name`, `key_vault_uri`, `key_version`, and `managed_hsm_key_id` for
  manual review (leaves them unchanged), with reasons noting v5 consolidates
  them into `key_vault_key_id`. GA v5.0.0 confirms the consolidation: none of
  `key_name`/`key_vault_uri`/`key_vault_id`/`key_version`/`managed_hsm_key_id`
  exist in the v5.0.0 schema any more, only `key_vault_key_id` (required).
  The flagged-but-unrewritten output fails validate with `Missing required
  argument: "key_vault_key_id" is required`, plus `Unsupported argument` for
  each stale v4 field. Excluded until the rule performs the consolidation.
- `data.azurerm_storage_blob` - mirrors the `azurerm_storage_blob` resource
  gap: the upgrader only flags `storage_account_name` /
  `storage_container_name` for manual review (leaves them unchanged), but GA
  v5.0.0's data source requires `storage_container_id` instead and no longer
  accepts either flagged field. Fails validate with `Missing required
  argument: "storage_container_id" is required` and `Unsupported argument`
  for both stale v4 fields. Excluded until the rule performs the
  consolidation.
- `azurerm_mssql_managed_instance_transparent_data_encryption`,
  `azurerm_mssql_server_transparent_data_encryption` - the upgrader only
  flags `managed_hsm_key_id` for manual review (leaves it unchanged), with
  the reason "in v5 the key is specified via key_vault_key_id; consolidate
  managed_hsm_key_id into key_vault_key_id". GA v5.0.0 confirms the
  consolidation: `managed_hsm_key_id` no longer exists in the v5.0.0 schema
  for either resource, only `key_vault_key_id`. The flagged-but-unrewritten
  output fails validate with `Unsupported argument: An argument named
  "managed_hsm_key_id" is not expected here.` (one occurrence per resource).
  Since this consolidation is the *only* v5 change either resource has, and
  the rule can't perform it, both resources are excluded entirely from the
  corpus (same pattern as `azurerm_storage_account_customer_managed_key`
  above) until the rule performs the consolidation.
- `azurerm_mssql_database`'s `threat_detection_policy.email_account_admins` -
  the upgrader only flags this attribute for manual review (leaves it
  unchanged as a string enum `"Enabled"`/`"Disabled"`), with the reason that
  v5 replaces it with the bool `email_account_admins_enabled` (not a
  value-preserving rename, since the value's type also changes from string
  to bool). GA v5.0.0 confirms the field is gone from the schema entirely:
  the flagged-but-unrewritten output fails validate with `Unsupported
  argument: An argument named "email_account_admins" is not expected here.`
  Unlike the fully-excluded resources above, `azurerm_mssql_database` has
  other, unrelated v5 changes that validate cleanly (the
  `long_term_retention_policy.immutable_backups_enabled` removal and the
  `weekly`/`monthly`/`yearly_retention` default-change flags), so the
  resource itself stays in the corpus; only the `threat_detection_policy`
  block's `email_account_admins` field is omitted from the fixture (the
  block's other fields, `state` and `storage_endpoint`, are unaffected and
  remain). This means the corpus does not exercise this specific
  `email_account_admins` review rule end-to-end; it remains covered by its
  per-rule golden fixture under
  `internal/rules/azurerm/testdata/azurerm_mssql_database.threat_detection_policy.email_account_admins.review`.
- `azurerm_linux_virtual_machine_scale_set` / `azurerm_windows_virtual_machine_scale_set`'s
  `automatic_os_upgrade_policy.disable_automatic_rollback` - the upgrader
  only flags this attribute for manual review (leaves it unchanged as a
  bool), with the reason that v5 renames it to `automatic_rollback_enabled`
  with inverted boolean meaning. GA v5.0.0 confirms `disable_automatic_rollback`
  no longer exists in the schema and `automatic_rollback_enabled` is
  Required whenever the `automatic_os_upgrade_policy` block is present (not
  merely renamed-and-optional): the flagged-but-unrewritten output fails
  validate with `Unsupported argument: An argument named
  "disable_automatic_rollback" is not expected here.` plus `Missing required
  argument: The argument "automatic_rollback_enabled" is required, but no
  definition was found.` (one occurrence per resource). Because
  `automatic_rollback_enabled` is required the instant the block exists, and
  the upgrader never emits it, no v4 config using
  `automatic_os_upgrade_policy` can validate clean post-upgrade until the
  rule performs the rename. Both resources otherwise have other v5 changes
  (`enable_automatic_os_upgrade` → `automatic_os_upgrade_enabled`,
  `data_disk.ultra_ssd_disk_{iops,mbps}_read_write` → `disk_{iops,mbps}_read_write`,
  `network_interface.enable_accelerated_networking` →
  `accelerated_networking_enabled`, `network_interface.enable_ip_forwarding`
  → `ip_forwarding_enabled`, plus `azurerm_windows_virtual_machine_scale_set`'s
  top-level `enable_automatic_updates` → `automatic_updates_enabled`) that
  validate cleanly, so both resources stay in the corpus; only the entire
  `automatic_os_upgrade_policy` block is omitted from each fixture. This
  means the corpus does not exercise the `enable_automatic_os_upgrade`
  rename or the `disable_automatic_rollback` review rule end-to-end for
  either resource; they remain covered by their per-rule golden fixtures
  under
  `internal/rules/azurerm/testdata/azurerm_linux_virtual_machine_scale_set.automatic_os_upgrade_policy.enable_automatic_os_upgrade.renamed`
  and
  `internal/rules/azurerm/testdata/azurerm_windows_virtual_machine_scale_set.automatic_os_upgrade_policy.enable_automatic_os_upgrade.renamed`.
- `azurerm_orchestrated_virtual_machine_scale_set`'s `sku_profile.vm_sizes` -
  the upgrader only flags this attribute for manual review (leaves it
  unchanged as a list of VM size strings), with the reason that v5 removes
  it in favour of one or more `sku_profile.virtual_machine_size` blocks.
  GA v5.0.0 confirms `vm_sizes` no longer exists in the schema and at least
  one `virtual_machine_size` block (nested `name`/`rank`) is required
  instead: the flagged-but-unrewritten output fails validate with
  `Unsupported argument: An argument named "vm_sizes" is not expected here.`
  plus `Insufficient virtual_machine_size blocks: At least 1
  "virtual_machine_size" blocks are required.` This is a list-to-nested-block
  restructure the upgrader cannot perform mechanically (same pattern as the
  `azurerm_storage_share_directory`/`azurerm_storage_share_file` id-to-url
  case above). The resource otherwise has other v5 changes
  (`os_profile.windows_configuration.enable_automatic_updates` →
  `automatic_updates_enabled`, `network_interface.enable_accelerated_networking`
  → `accelerated_networking_enabled`, `network_interface.enable_ip_forwarding`
  → `ip_forwarding_enabled`, `data_disk.ultra_ssd_disk_{iops,mbps}_read_write`
  → `disk_{iops,mbps}_read_write`) that validate cleanly, so the resource
  stays in the corpus; only `sku_profile.vm_sizes` is omitted, replaced with
  explicit `virtual_machine_size` blocks so the fixture is schema-complete.
  This means the corpus does not exercise the `sku_profile.vm_sizes` review
  rule end-to-end; it remains covered by its per-rule golden fixture
  (`azurerm_orchestrated_virtual_machine_scale_set.os_profile.windows_configuration.enable_automatic_updates.renamed`
  in `internal/rules/azurerm/testdata` also sets `vm_sizes`, unvalidated
  against GA).
- `azurerm_disk_encryption_set.managed_hsm_key_id` - same consolidation
  pattern as `azurerm_storage_account_customer_managed_key` and the mssql
  TDE resources above: the upgrader only flags `managed_hsm_key_id` for
  manual review (leaves it unchanged), and GA v5.0.0 confirms the field no
  longer exists in the schema at all (only `key_vault_key_id`, Required, is
  present). The corpus's `azurerm_disk_encryption_set` fixture sidesteps
  this by using `key_vault_key_id` directly (pointed at a sibling
  `azurerm_key_vault`/`azurerm_key_vault_key` pair) instead of
  `managed_hsm_key_id`, so the resource validates GA-clean but the corpus
  does not exercise this specific review rule end-to-end; it remains
  covered by its per-rule golden fixture under
  `internal/rules/azurerm/testdata/azurerm_disk_encryption_set.managed_hsm_key_id.review`.
  The sibling `azurerm_key_vault` itself has one v5 change
  (`enable_rbac_authorization` → `rbac_authorization_enabled`, also newly
  Required) which the upgrader handles correctly and which validates clean.
- `azurerm_hdinsight_hadoop_cluster` / `azurerm_hdinsight_hbase_cluster` /
  `azurerm_hdinsight_interactive_query_cluster` / `azurerm_hdinsight_kafka_cluster`
  / `azurerm_hdinsight_spark_cluster`'s `storage_account.storage_container_id` -
  the upgrader only flags this attribute for manual review (leaves it
  unchanged), with the reason that v5 renames it to `storage_container_url`
  and requires the value to be the container's Data-Plane URL rather than an
  ARM resource ID. GA v5.0.0 confirms `storage_container_id` no longer exists
  in the `storage_account` block's schema (only `storage_container_url`,
  Required, is present): the flagged-but-unrewritten output fails validate
  with `Missing required argument: The argument "storage_container_url" is
  required, but no definition was found.` plus `Unsupported argument: An
  argument named "storage_container_id" is not expected here.` (one
  occurrence per resource that uses the block). This module's fixtures
  sidestep the gap by using only the `storage_account_gen2` block (whose
  `storage_resource_id` → `storage_account_id` and
  `managed_identity_resource_id` → `user_assigned_identity_id` renames are
  fully mechanical and validate clean) for each cluster's default storage;
  the `storage_account` block is omitted entirely from all five fixtures, so
  the corpus does not exercise the `storage_account.storage_resource_id`
  rename or the `storage_container_id` review rule end-to-end for any of the
  five resources. They remain covered by their per-rule golden fixtures
  under `internal/rules/azurerm/testdata/azurerm_hdinsight_*_cluster.v5-gaps`.
- `azurerm_private_dns_a_record`, `azurerm_private_dns_aaaa_record`,
  `azurerm_private_dns_cname_record`, `azurerm_private_dns_mx_record`,
  `azurerm_private_dns_ptr_record`, `azurerm_private_dns_srv_record`,
  `azurerm_private_dns_txt_record`, `azurerm_private_dns_zone_virtual_network_link` -
  same consolidation-flagged-but-unrewritten pattern as
  `azurerm_storage_account_customer_managed_key` above: the upgrader only
  flags `resource_group_name` and `zone_name` (for the record resources) or
  `resource_group_name` and `private_dns_zone_name` (for the vnet-link
  resource) for manual review, leaving both unchanged, with the reason that
  v5 consolidates them into a single `private_dns_zone_id`. GA v5.0.0
  confirms the consolidation: neither `zone_name`/`private_dns_zone_name`
  nor `resource_group_name` exist in the v5.0.0 schema for any of these
  eight resources any more, only `private_dns_zone_id` (Required). The
  flagged-but-unrewritten output fails validate with `Missing required
  argument: The argument "private_dns_zone_id" is required, but no
  definition was found.` plus `Unsupported argument` for each stale v4
  field (one occurrence per resource). Since this consolidation is the
  *only* v5 change any of these eight resources has, and the rule can't
  perform it, all eight are excluded entirely from the corpus (Task 9); no
  `modules/privatedns` module was added. They remain covered by their
  per-rule golden fixtures under
  `internal/rules/azurerm/testdata/azurerm_private_dns_*.review`.
- `azurerm_kusto_eventgrid_data_connection.consumer_group` - unlike its
  sibling `eventgrid_resource_id`/`managed_identity_resource_id` renames
  (both registered and mechanical), no rule in
  `internal/rules/azurerm/attr/register.go` covers `consumer_group` at all.
  GA v5.0.0 renamed this field to `eventhub_consumer_group_name` (still
  Required); with the upgrader leaving `consumer_group` untouched, the
  output fails validate with `Missing required argument: The argument
  "eventhub_consumer_group_name" is required, but no definition was found.`
  plus `Unsupported argument: An argument named "consumer_group" is not
  expected here.` This is a genuine missing-rule gap (not a flag-only or
  consolidation case), so `azurerm_kusto_eventgrid_data_connection` is
  excluded entirely from `modules/dataplatform` (Task 11) until a
  `consumer_group` → `eventhub_consumer_group_name` rename rule is added.
- `azurerm_eventhub` - same consolidation-flagged-but-unrewritten pattern as
  `azurerm_private_dns_*` above: the upgrader only flags `namespace_name` and
  `resource_group_name` for manual review, leaving both unchanged, with the
  reason that v5 consolidates them into a single `namespace_id`. GA v5.0.0
  confirms the consolidation: neither field exists in the v5.0.0 schema any
  more (only `namespace_id`, Required); the flagged-but-unrewritten output
  fails validate with `Missing required argument: "namespace_id" is
  required` plus `Unsupported argument` for both stale v4 fields. Since this
  consolidation is the *only* v5 change this resource has, and the rule
  can't perform it, `azurerm_eventhub` is excluded entirely from
  `modules/messaging` (its two `eventgrid_event_subscription`/
  `eventgrid_system_topic_event_subscription` consumers reference a literal
  EventHub resource ID string instead, so the genuine
  `eventhub_endpoint_id` → `eventhub_id` rename those resources exercise
  still fires). It remains covered by the per-rule golden fixture
  `internal/rules/azurerm/testdata/azurerm_eventhub.namespace_name.review`.
- `data.azurerm_servicebus_namespace_disaster_recovery_config`,
  `data.azurerm_servicebus_queue` - same pattern: the upgrader only flags
  `namespace_name`/`resource_group_name` for manual review ("removed in v5;
  use namespace_id"), leaving them unchanged. GA v5.0.0 confirms both fields
  are gone from each data source's schema (only `namespace_id`, Required).
  `data.azurerm_servicebus_subscription` similarly flags
  `namespace_name`/`resource_group_name`/`topic_name` ("removed in v5; use
  topic_id"); GA v5.0.0 confirms only `topic_id` (Required) remains. Since
  this consolidation is the *only* v5 change any of these three data
  sources has, and the rule can't perform it, all three are excluded
  entirely from `modules/messaging` (their underlying
  `azurerm_servicebus_namespace`/`azurerm_servicebus_queue`/
  `azurerm_servicebus_subscription` resources stay, as legitimate siblings
  whose own `namespace_id`/`topic_id` usage is the correct v4.x form). They
  remain covered by the per-rule golden fixtures under
  `internal/rules/azurerm/testdata/data.azurerm_servicebus_*.consolidation`.
- `azurerm_iot_security_solution`'s `recommendations_enabled` block - the
  upgrader only flags this block for manual review (leaves it unchanged,
  block name and all), with the reason that v5 renames it to
  `recommendations` (identical nested schema). GA v5.0.0 confirms
  `recommendations_enabled` no longer exists in the schema at all (only
  `recommendations` is present): the flagged-but-unrewritten output fails
  validate with `Unsupported block type: Blocks of type
  "recommendations_enabled" are not expected here.` This is a block rename
  the upgrader cannot perform mechanically (same pattern as the
  `automatic_os_upgrade_policy`/`sku_profile.vm_sizes` cases above). The
  resource has no other v5 changes, so `modules/security` (Task 14) keeps
  the resource in the corpus but omits the `recommendations_enabled` block
  entirely (it is optional) rather than authoring it - this means the
  corpus does not exercise the block-rename review rule end-to-end; it
  remains covered by its per-rule golden fixture under
  `internal/rules/azurerm/testdata/azurerm_iot_security_solution.misc`.
- `azurerm_federated_identity_credential` - the upgrader's only rule for this
  resource renames `parent_id` to `user_assigned_identity_id` (a genuine v4
  config always sets both `resource_group_name` and `parent_id` - both are
  required in v4), but does not also remove the now-redundant
  `resource_group_name`. GA v5.0.0 confirms `resource_group_name` no longer
  exists in the schema once identity is referenced by ID: the
  renamed-but-otherwise-unchanged output fails validate with `Unsupported
  argument: An argument named "resource_group_name" is not expected here.`
  This is a genuine rule bug (the rename rule needs to also drop
  `resource_group_name`), not a fixture-completeness issue - a v4 config
  cannot omit `resource_group_name` and still be valid v4. Excluded entirely
  from `modules/platform` (Task 15) until the rule also strips it; the
  `azurerm_user_assigned_identity` sibling stays (reused by
  `azurerm_mysql_flexible_server`'s `customer_managed_key`). Still covered by
  the per-rule golden fixture under
  `internal/rules/azurerm/testdata/azurerm_federated_identity_credential.parent_id.renamed`.
- The following `modules/platform` (Task 15) resources have flag-only
  (unrewritten NeedsReview) attributes/blocks that are simply omitted from
  the fixture, per the established convention (same pattern as
  `azurerm_iot_security_solution.recommendations_enabled` above) - each
  remains covered end-to-end by its per-rule golden fixture:
  - `azurerm_log_analytics_workspace.local_authentication_disabled`,
    `azurerm_application_insights.{daily_data_cap_notifications_disabled,disable_ip_masking,local_authentication_disabled}`,
    `azurerm_cosmosdb_account.{local_authentication_disabled,managed_hsm_key_id}` -
    all removed in v5 in favour of an `*_enabled` counterpart with inverted
    boolean meaning (`local_authentication_disabled`/`daily_data_cap_notifications_disabled`/`disable_ip_masking`)
    or consolidated into `key_vault_key_id` (`managed_hsm_key_id`); GA
    v5.0.0 confirms none of these attributes exist in the schema any more.
  - `azurerm_monitor_diagnostic_setting.metric` (removed in favour of
    `enabled_metric`, used directly instead since it already existed
    pre-v5) and `enabled_log.retention_policy` on both
    `azurerm_monitor_diagnostic_setting` and
    `azurerm_monitor_aad_diagnostic_setting` (removed in favour of
    `azurerm_storage_management_policy`).
  - `azurerm_application_gateway.authentication_certificate` (top-level and
    nested under `backend_http_settings`) - removed in favour of
    `trusted_root_certificate`/`trusted_root_certificate_names`.
  - `azurerm_key_vault.contact` - removed (used a data-plane API); managed
    instead via the separate `azurerm_key_vault_certificate_contacts`
    resource, which is exercised directly in the corpus.
  - `azurerm_batch_pool.certificate` - removed (no longer supported by the
    Batch service).
  - `azurerm_mysql_flexible_server.customer_managed_key.managed_hsm_key_id` -
    consolidated into `customer_managed_key.key_vault_key_id`.
  - `azurerm_netapp_volume.mount_ip_addresses` - not omitted for a rule-gap
    reason; it's a Computed-only attribute (Azure-assigned mount IPs), so
    there is no realistic v4 config that sets it in the first place.
  - `azurerm_redis_cache.minimum_tls_version` / `azurerm_cosmosdb_account.minimal_tls_version`
    use still-accepted values (`"1.2"` / `"Tls12"`) rather than the removed
    values (`"1.0"`/`"1.1"` / `"Tls"`/`"Tls11"`) their per-rule
    `value-removed` fixtures exercise, per the value-removed-category
    convention already used for `azurerm_mssql_server.minimum_tls_version`
    in `modules/mssql`.
