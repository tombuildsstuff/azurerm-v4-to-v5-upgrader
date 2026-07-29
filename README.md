# `tombuildsstuff/azurerm-v4-to-v5-upgrader`

This tool is a CLI application which upgrades a Terraform configuration from using v4.x of the `hashicorp/terraform-provider-azurerm`
to using v5.0.0 - applying the breaking changes from the [official 5.0 upgrade guide](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/5.0-upgrade-guide) directly to your `.tf` files.

It does this by parsing the HCL, rewriting the constructs which have changes in v5 (renaming attributes/blocks,
updating `name` -> `id` references, removing resources, value changes and other changes flagged in the
upgrade guide) as needed for v5.0. It does this whilst preserving comments and interpolations.

This intentionally **does not** commit these changes, so that you can easily diff them as needed and
covers **329 breaking-changes** across Resources, Data Sources and the Provider block.

Where changes can be made automatically, the tool makes them - where they can't, it flags them in
the output for your attention. This is built both from the upgrade guide itself, and a diff of v4.81.0 ->
v5.0.0 to confirm there's no missing breaking changes).

It's recommended to upgrade to the latest version of 4.x.x prior to running this tool, but you can determine
if there's any pending changes to be made by running  `./azurerm-v4-to-v5-upgrader validate`. To run it
for real, making those changes, you can run `./azurerm-v4-to-v5-upgrader upgrade`.

## Install

Requires Go 1.26+.

```sh
go install github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader@latest
```

This installs the `upgrader` binary into `$(go env GOPATH)/bin`. Or build from a
checkout:

```sh
go build -o upgrader .
```

## Usage

The tool has two commands: `validate` (read-only) and `update` (rewrites files).

### `validate` - check whether an upgrade is needed

Recursively scans a directory tree and reports what would change, **without
writing anything**. Exits `1` if any configuration needs an upgrade, `0` if not -
so it works as a CI gate.

```sh
upgrader validate            # scan the current directory tree
upgrader validate ./modules  # scan a specific path
upgrader validate --report json
```

### `update` - apply the upgrade

Rewrites `.tf` files in place. By default it recurses into nested module
directories (matching `validate`).

```sh
upgrader update --dir .              # upgrade the current tree in place
upgrader update --dir . --dry-run    # print a diff instead of writing
upgrader update --dir . --report json
```

Preview first with `--dry-run`, and make sure your work is committed to version
control before running without it.

#### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--dir` | `.` | Module directory to upgrade. |
| `--dry-run` | `false` | Print a unified diff instead of writing files. |
| `--recursive` | `true` | Recurse into nested module directories. Use `--recursive=false` for the top-level directory only. |
| `--fmt` | `true` | Canonically format changed files, like `terraform fmt`. Use `--fmt=false` to keep the original formatting on lines the upgrade did not change (see [Output formatting](#output-formatting)). |
| `--report` | `text` | Report format: `text` or `json`. |

Directories named `.terraform` are always skipped.

## What it does

Findings come in three severities:

- **Changed** - a value-preserving rewrite the tool applied automatically
  (attribute/block renames, `name` → `id` reference rewrites, `bool` → enum
  conversions, list → repeated-block conversions, field moves into nested
  blocks).
- **Note** - a default value that changed in v5, pinned to keep v4 behaviour.
- **Needs review** - a breaking change with no single safe mechanical rewrite,
  surfaced for you to resolve by hand (removed resources/data sources, removed or
  now-required attributes/blocks, dropped enum values).

A run reports how many `azurerm_` resources it saw and how many rules fired, e.g.:

```
azurerm v5 upgrade: 52 resource(s) seen, 7 rule(s) fired
```

## Output formatting

To apply an edit the tool re-serialises each changed file through `hclwrite`,
which canonically formats the whole file - the same formatting `terraform fmt`
applies. So a changed file can show incidental formatting-only diffs on lines the
upgrade never touched (most visibly, the JSON-style object syntax `"k": v` is
rewritten to `"k" : v`, exactly as `terraform fmt` does). Files with no transform
are left byte-for-byte unchanged.

Pass `update --fmt=false` to opt out: the tool then restores the original text on
any line it changed only in whitespace, so the diff contains just the real
upgrade edits. Alternatively, run `terraform fmt` before the upgrader if you
prefer the module already canonical.

## GitHub Actions

`validate --github-actions` additionally emits GitHub Actions channels: inline
annotations, a step-summary section, an action output, and - on pull requests - a
summary comment. It reads the standard `GITHUB_*` environment variables
(`GITHUB_TOKEN`, `GITHUB_REPOSITORY`, `GITHUB_EVENT_NAME`, `GITHUB_EVENT_PATH`,
`GITHUB_STEP_SUMMARY`, `GITHUB_OUTPUT`, `GITHUB_API_URL`) and sets the
`needs-changes` output.

```yaml
- name: Check azurerm v5 upgrade
  run: upgrader validate --github-actions
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

The step fails (exit `1`) when any configuration needs an upgrade.

## Testing

The tool is covered by a golden-fixture harness
(`internal/rules/rulestest`): each rule has `input.tf` / `expected.tf` /
`expected-report.json` fixtures, every registered rule ID must fire in at least
one fixture, and the harness parse-validates every transformed output so no rule
can emit invalid HCL. Engine integration tests run a full corpus end-to-end
(including idempotency and comment/interpolation preservation).

```sh
go test ./...
```

An opt-in acceptance layer runs `terraform validate` against the real azurerm v5
provider schema to prove the output is accepted (not just well-formed). It is
gated behind `AZURERM_ACC=1` and needs `terraform` plus a v5 provider available:

```sh
AZURERM_ACC=1 go test ./internal/engine/ -run TestAcceptanceV5Validate -v
```

## Development

The `Makefile` wraps the common development workflows:

| Target | What it does |
|--------|--------------|
| `make build` | Tidies the module, refreshes the (git-ignored) `vendor/` tree, and compiles the `upgrader` binary into `./upgrader`. |
| `make install` | Runs `build`, then installs the binary into `$(go env GOPATH)/bin`. |
| `make test` | Runs the full Go test suite (`go test ./...`). The opt-in `terraform validate` acceptance layer stays skipped unless `AZURERM_ACC=1` is set. |
| `make fmt` | Formats Go sources (`gofmt` + `goimports`, with imports grouped by module) and Terraform fixtures (`terraform fmt -recursive`). |
| `make tools` | Installs the developer dependencies the other targets use (currently `goimports`; `terraform` is assumed to already be on `PATH`). |

## License

Licensed under the Apache License, Version 2.0. See [`LICENSE`](LICENSE).
