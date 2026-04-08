# Provenance Collector

A Kubernetes-native tool that discovers running container images and Helm
releases, resolves digests, verifies signatures, checks for available updates,
and generates compliance-grade provenance reports.

Deployed as a **Kubernetes CronJob** and packaged as a
[Nebari Software Pack](https://github.com/nebari-dev/nebari-operator).

## What It Does

| Capability | Description |
|---|---|
| **Image Discovery** | Scans all pods across namespaces, deduplicates by workload owner |
| **Digest Resolution** | Resolves every image tag to its immutable SHA256 digest |
| **Signature Verification** | Checks for cosign signatures (existence or key-based verification) |
| **SBOM Detection** | Detects attached SPDX / CycloneDX attestations |
| **Update Checking** | Compares running tags against latest available semver tags |
| **Helm Release Tracking** | Discovers all deployed Helm releases with chart versions |
| **Provenance Reports** | Outputs timestamped JSON reports to PVC or ConfigMap |

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

## Configuration

All configuration is via environment variables, set through `values.yaml`:

| Variable | Default | Description |
|---|---|---|
| `PROVENANCE_NAMESPACES` | *(all)* | Comma-separated namespaces to scan |
| `PROVENANCE_EXCLUDE_NAMESPACES` | *(none)* | Namespaces to skip |
| `PROVENANCE_VERIFY_SIGNATURES` | `true` | Check cosign signatures |
| `PROVENANCE_COSIGN_PUBLIC_KEY` | *(empty)* | Path/KMS URI for cosign key |
| `PROVENANCE_HELM_ENABLED` | `true` | Discover Helm releases |
| `PROVENANCE_CHECK_UPDATES` | `true` | Check for newer image tags |
| `PROVENANCE_REPORT_OUTPUT` | `pvc` | Report output: `pvc` or `configmap` |
| `PROVENANCE_REPORT_PATH` | `/reports` | Filesystem path for PVC reports |
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
  helmEnabled: true
  checkUpdates: true
  reportOutput: "pvc"         # "pvc" or "configmap"

persistence:
  enabled: true
  size: 1Gi

# Nebari integration (optional)
nebariapp:
  enabled: false
```

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
cmd/provenance-collector/     Entry point
internal/
  config/                     Environment-based configuration
  kubernetes/                 Client factory (in-cluster + kubeconfig)
  discovery/
    images.go                 Pod-based container image discovery
    helm.go                   Helm release discovery via Helm SDK
  registry/
    digest.go                 Digest resolution via go-containerregistry
    updates.go                Semver-based update checking
  verify/
    cosign.go                 Signature verification via cosign
    sbom.go                   SBOM attestation detection
  report/
    types.go                  Report JSON schema types
    generator.go              Orchestrator with concurrent enrichment
    writer.go                 PVC and ConfigMap output writers
chart/                        Helm chart (CronJob + RBAC + NebariApp)
```

## License

BSD-3-Clause — see [LICENSE](LICENSE).
