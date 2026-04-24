<p align="center">
  <a href="https://nebari.dev">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/nebari-dev/nebari-design/main/logo-mark/horizontal/standard/Nebari-Logo-Horizontal-Lockup-White-text.png">
      <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/nebari-dev/nebari-design/main/logo-mark/horizontal/standard/Nebari-Logo-Horizontal-Lockup.png">
      <img alt="Nebari" src="https://raw.githubusercontent.com/nebari-dev/nebari-design/main/logo-mark/horizontal/standard/Nebari-Logo-Horizontal-Lockup.png" width="300">
    </picture>
  </a>
</p>

# Provenance Collector

A Kubernetes-native tool that discovers running container images and Helm
releases, resolves digests, verifies signatures, detects SLSA provenance
and SBOM attestations, checks for available updates, and generates
compliance-grade provenance reports.

Deployed as a **Kubernetes CronJob** and packaged as a
[Nebari Software Pack](https://github.com/nebari-dev/nebari-operator).

## What It Does

| Capability | Description |
|---|---|
| **Image Discovery** | Scans all pods across namespaces, deduplicates by workload owner |
| **Digest Resolution** | Resolves every image tag to its immutable SHA256 digest |
| **Signature Verification** | Checks for cosign signatures (existence or key-based verification) |
| **SLSA Provenance** | Detects SLSA provenance attestations via OCI referrers API |
| **SBOM Detection** | Detects attached SPDX / CycloneDX attestations |
| **Update Checking** | Compares running tags against latest semver tags (configurable level, pre-release filtering) |
| **Helm Release Tracking** | Discovers all deployed Helm releases with chart versions |
| **Web Dashboard** | Optional UI with filters, sorting, pagination, and image detail panel |
| **Grafana Integration** | JSON API compatible with the Infinity datasource for dashboards and alerting |
| **Provenance Reports** | Outputs timestamped JSON reports to PVC or ConfigMap with automatic retention |

## Quick Start

### Prerequisites

- Kubernetes cluster (1.26+)
- Helm 3
- `kubectl` configured for your cluster

### Install

```bash
helm install provenance-collector chart/ \
  --namespace provenance-system \
  --create-namespace
```

### Trigger a Manual Run

```bash
kubectl create job --from=cronjob/provenance-collector \
  manual-run -n provenance-system
```

### View the Report

```bash
# PVC-based (default)
kubectl logs -n provenance-system job/manual-run

# ConfigMap-based
kubectl get configmap provenance-report \
  -n provenance-system \
  -o jsonpath='{.data.report\.json}' | jq .
```

## Web Dashboard

The collector includes an optional web UI for browsing provenance reports.

```yaml
webUI:
  enabled: true
```

Access via port-forward:

```bash
kubectl port-forward svc/provenance-collector-web 8080:8080 -n provenance-system
# Open http://localhost:8080
```

The dashboard provides:

- Summary stat cards (signed %, SLSA provenance, SBOM, updates available)
- Report timeline to browse historical reports
- Filterable, sortable, paginated image table
- Click any image row for a detail panel showing signature, SLSA, SBOM, and update info
- Helm releases table

When deployed with `nebariapp.enabled: true`, the dashboard is accessible through the
Nebari gateway with OIDC authentication.

### Dashboard API

The web dashboard exposes a JSON API that can be used by external tools:

| Endpoint | Description |
|---|---|
| `GET /api/reports` | List all reports (newest first) with summary |
| `GET /api/reports/latest` | Get the most recent report |
| `GET /api/reports/<filename>` | Get a specific report by filename |
| `GET /healthz` | Health check |

## Grafana Integration

The provenance data can be surfaced in Grafana using the
[Infinity datasource](https://grafana.com/grafana/plugins/yesoreyeram-infinity-datasource/)
plugin, which queries the dashboard's JSON API.

### Setup

1. Enable the web dashboard (`webUI.enabled: true`)
2. Install the Infinity datasource in Grafana
3. Add a datasource pointing at the dashboard service:
   - **URL:** `http://provenance-collector-web.provenance-system.svc:8080`
   - **Type:** JSON

### Example panels

**Stat panel** (unique images count):
- Type: JSON, URL: `/api/reports/latest`
- Column: `summary.uniqueImages`

**Images table**:
- Type: JSON, URL: `/api/reports/latest`, Root: `images`
- Columns: `image`, `namespace`, `signature.signed`, `provenance.hasProvenance`, `update.updateAvailable`

**Alerting** (images with updates):
```
WHEN count() OF images WHERE updateAvailable = true IS ABOVE 0
```

An example dashboard (11 panels covering unique images, signature status, SLSA provenance, Helm releases, and more) is available at [`examples/grafana-dashboard.json`](examples/grafana-dashboard.json). Import it directly into Grafana as a `dashboard.grafana.app/v2beta1` resource.

See [docs/configuration.md](docs/configuration.md) for the full Grafana setup guide.

## Configuration

All configuration is via environment variables, set through `values.yaml`:

| Variable | Default | Description |
|---|---|---|
| `PROVENANCE_NAMESPACES` | *(all)* | Comma-separated namespaces to scan |
| `PROVENANCE_EXCLUDE_NAMESPACES` | *(none)* | Namespaces to skip |
| `PROVENANCE_VERIFY_SIGNATURES` | `true` | Check cosign signatures |
| `PROVENANCE_COSIGN_PUBLIC_KEY` | *(empty)* | Path/KMS URI for cosign key |
| `PROVENANCE_CHECK_SBOM` | `true` | Check for SBOM attestations |
| `PROVENANCE_CHECK_PROVENANCE` | `true` | Check for SLSA provenance attestations |
| `PROVENANCE_HELM_ENABLED` | `true` | Discover Helm releases |
| `PROVENANCE_CHECK_UPDATES` | `true` | Check for newer image tags |
| `PROVENANCE_UPDATE_LEVEL` | `patch` | Min version bump to flag: `patch`, `minor`, or `major` |
| `PROVENANCE_SKIP_PRERELEASE` | `true` | Ignore alpha/beta/RC versions in updates |
| `PROVENANCE_REPORT_OUTPUT` | `pvc` | Report output: `pvc` or `configmap` |
| `PROVENANCE_REPORT_PATH` | `/reports` | Filesystem path for PVC reports |
| `PROVENANCE_REPORT_RETENTION` | `168h` | Auto-prune reports older than this (`-1` to disable) |
| `PROVENANCE_REGISTRY_TIMEOUT` | `30s` | Timeout for registry operations |
| `PROVENANCE_CLUSTER_NAME` | *(empty)* | Cluster name in report metadata |

See [docs/configuration.md](docs/configuration.md) for the full reference.

## Report Format

Reports are JSON with this structure:

```json
{
  "metadata": {
    "generatedAt": "2025-01-15T06:00:00Z",
    "collectorVersion": "0.1.0",
    "clusterName": "production",
    "namespacesScanned": ["default", "monitoring"]
  },
  "images": [
    {
      "image": "nginx:1.27-alpine",
      "digest": "sha256:abc123...",
      "namespace": "default",
      "workload": { "kind": "Deployment", "name": "nginx" },
      "signature": { "signed": true, "verified": true },
      "provenance": { "hasProvenance": true, "predicateType": "https://slsa.dev/provenance/v1" },
      "sbom": { "hasSBOM": true, "format": "spdx" },
      "update": {
        "currentTag": "1.27",
        "latestInMajor": "1.27.3",
        "updateAvailable": true
      }
    }
  ],
  "helmReleases": [
    {
      "releaseName": "ingress-nginx",
      "namespace": "ingress",
      "chart": "ingress-nginx",
      "version": "4.8.0",
      "appVersion": "1.9.4",
      "status": "deployed"
    }
  ],
  "summary": {
    "totalImages": 42,
    "uniqueImages": 28,
    "signedImages": 15,
    "verifiedImages": 12,
    "imagesWithSBOM": 10,
    "imagesWithProvenance": 3,
    "imagesWithUpdates": 5,
    "totalHelmReleases": 8,
    "helmReleasesWithUpdates": 2
  }
}
```

See [docs/report-schema.md](docs/report-schema.md) for the full schema reference.

## Helm Chart Values

Key values (see `chart/values.yaml` for all options):

```yaml
schedule: "0 6 * * *"        # Daily at 6 AM UTC
config:
  namespaces: []              # Empty = all namespaces
  excludeNamespaces: []
  verifySignatures: true
  checkSBOM: true
  checkProvenance: true
  checkUpdates: true
  updateLevel: "patch"        # "patch", "minor", or "major"
  skipPrerelease: true
  reportOutput: "pvc"         # "pvc" or "configmap"
  reportRetention: "168h"     # 1 week, "-1" to keep forever

persistence:
  enabled: true
  size: 1Gi

webUI:
  enabled: false              # Enable the web dashboard

# Nebari integration (optional)
nebariapp:
  enabled: false
```

See [examples/](examples/) for complete deployment examples (standalone, Nebari, ArgoCD).

## RBAC

The collector requires cluster-wide read access:

| Resource | Verbs | Purpose |
|---|---|---|
| `pods`, `namespaces` | get, list | Image discovery |
| `deployments`, `replicasets`, `statefulsets`, `daemonsets` | get, list | Owner resolution |
| `jobs`, `cronjobs` | get, list | Owner resolution |
| `secrets` | get, list | Helm release storage |
| `configmaps` | get, list, create, update | Report output |

## Development

```bash
# Build
make build

# Test
make test

# Lint
make lint

# Docker
make docker-build
```

### Local Testing with kind

```bash
kind create cluster
docker build -t provenance-collector:dev .
kind load docker-image provenance-collector:dev

helm install pc chart/ \
  --set image.repository=provenance-collector \
  --set image.tag=dev \
  --set image.pullPolicy=Never \
  --set config.reportOutput=configmap \
  --set persistence.enabled=false \
  --set config.verifySignatures=false

kubectl create job --from=cronjob/pc-provenance-collector test-run
kubectl logs job/test-run
```

## Architecture

```
cmd/
  provenance-collector/       Collector entry point (CronJob)
  dashboard/                  Web dashboard entry point
internal/
  config/                     Environment-based configuration
  kubernetes/                 Client factory (in-cluster + kubeconfig)
  dashboard/                  HTTP server, API handlers, HTML UI
  discovery/
    images.go                 Pod-based container image discovery
    helm.go                   Helm release discovery via Helm SDK
  registry/
    digest.go                 Digest resolution via go-containerregistry
    updates.go                Semver-based update checking
  verify/
    cosign.go                 Signature verification via cosign
    sbom.go                   SBOM attestation detection
    provenance.go             SLSA provenance detection via OCI referrers
  report/
    types.go                  Report JSON schema types
    generator.go              Orchestrator with concurrent enrichment
    writer.go                 PVC and ConfigMap output writers
chart/                        Helm chart (CronJob + RBAC + Dashboard + NebariApp)
examples/                     Deployment examples (standalone, Nebari, ArgoCD)
```

## License

BSD-3-Clause — see [LICENSE](LICENSE).
