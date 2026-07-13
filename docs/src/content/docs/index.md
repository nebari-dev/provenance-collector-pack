---
title: Provenance Collector
description: Compliance-grade supply-chain provenance for every container image and Helm release running on a Nebari cluster.
---

The Provenance Collector is a **Nebari Software Pack** that produces
compliance-grade supply-chain reports for every container image and Helm
release running on a Kubernetes cluster. It is deployed by the
[Nebari Operator](https://github.com/nebari-dev/nebari-operator) as part of
Nebari Infrastructure Core (NIC), runs on a schedule as a `CronJob`, and ships
each timestamped JSON report to a web dashboard, a shared PVC, or a ConfigMap.

It exists because answering *"what is actually running on this cluster, where
did it come from, and is it signed?"* should not require manual auditing.

## What it does

| Capability | Description |
|---|---|
| **Image discovery** | Scans all pods across namespaces, deduplicates by workload owner |
| **Digest resolution** | Resolves every image tag to its immutable SHA256 digest |
| **Signature verification** | Checks for cosign signatures (existence or key-based verification) |
| **SLSA provenance** | Detects SLSA provenance attestations via the OCI referrers API |
| **SBOM detection** | Detects attached SPDX / CycloneDX attestations |
| **Update checking** | Compares running tags against the latest semver tags |
| **Helm release tracking** | Discovers all deployed Helm releases with chart versions |
| **Web dashboard** | Optional UI with filters, sorting, and an image detail panel |
| **Grafana integration** | JSON API compatible with the Infinity datasource |

## Reference

- [Configuration](/configuration/) — every environment variable and its chart value.
- [Report Schema](/report-schema/) — the JSON output structure.
- [NebariApp CRD](/nebariapp-crd-reference/) — operator integration fields.

> **Status:** Under active development as part of NIC. APIs, chart values, and
> report schema may change without notice while pre-1.0.

Full installation and deployment instructions live in the
[project README](https://github.com/nebari-dev/provenance-collector-pack#readme)
and the [`examples/`](https://github.com/nebari-dev/provenance-collector-pack/tree/main/examples)
directory. These will be migrated into this site in follow-up changes.
