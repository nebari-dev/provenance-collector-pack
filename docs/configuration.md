# Configuration Reference

The provenance collector is configured entirely via environment variables.
When deployed with the Helm chart, these are set through `values.yaml` and
injected via a ConfigMap.

## Environment Variables

### Namespace Filtering

| Variable | Type | Default | Description |
|---|---|---|---|
| `PROVENANCE_NAMESPACES` | comma-separated | *(all)* | Namespaces to scan. Empty means all namespaces. |
| `PROVENANCE_EXCLUDE_NAMESPACES` | comma-separated | *(none)* | Namespaces to exclude from scanning. |

### Registry

| Variable | Type | Default | Description |
|---|---|---|---|
| `PROVENANCE_REGISTRY_AUTH` | path | *(empty)* | Path to Docker `config.json` for private registry authentication. |
| `PROVENANCE_REGISTRY_TIMEOUT` | duration | `30s` | Timeout for registry operations (digest resolution, tag listing). |
| `PROVENANCE_CHECK_UPDATES` | bool | `true` | Check registries for newer semver tags. |

### Signature Verification

| Variable | Type | Default | Description |
|---|---|---|---|
| `PROVENANCE_VERIFY_SIGNATURES` | bool | `true` | Check for cosign signatures on images. |
| `PROVENANCE_COSIGN_PUBLIC_KEY` | string | *(empty)* | Path or KMS URI for cosign public key. Empty = existence check only (no trust chain verification). |

### Helm

| Variable | Type | Default | Description |
|---|---|---|---|
| `PROVENANCE_HELM_ENABLED` | bool | `true` | Discover deployed Helm releases. |

### Report Output

| Variable | Type | Default | Description |
|---|---|---|---|
| `PROVENANCE_REPORT_OUTPUT` | string | `pvc` | Output type: `pvc` (filesystem) or `configmap`. |
| `PROVENANCE_REPORT_PATH` | path | `/reports` | Directory for PVC-based reports. |
| `PROVENANCE_REPORT_CONFIGMAP` | string | `provenance-report` | ConfigMap name for configmap-based output. |
| `PROVENANCE_REPORT_CONFIGMAP_NAMESPACE` | string | *(release namespace)* | Namespace for ConfigMap output. |

### Metadata

| Variable | Type | Default | Description |
|---|---|---|---|
| `PROVENANCE_CLUSTER_NAME` | string | *(empty)* | Cluster name included in report metadata. Informational only. |

### Kubernetes

| Variable | Type | Default | Description |
|---|---|---|---|
| `KUBECONFIG` | path | *(empty)* | Path to kubeconfig file. Empty = use in-cluster credentials. |

## Duration Format

Duration values use Go's `time.ParseDuration` format:
- `30s` — 30 seconds
- `5m` — 5 minutes
- `1h30m` — 1 hour 30 minutes

## Boolean Values

Boolean values accept: `true`, `1` (truthy) or anything else (falsy).
