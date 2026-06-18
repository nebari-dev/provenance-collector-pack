+++
title = 'Provenance Collector'
description = "A Kubernetes-native CronJob that discovers running images and Helm releases, resolves digests, verifies signatures, detects SLSA provenance and SBOM attestations, checks for updates, and emits a timestamped JSON report. Optionally surfaced via a web dashboard and Grafana."
+++

The Provenance Collector is a **Nebari Software Pack** that produces compliance-grade
supply-chain reports for every container image and Helm release running on a Kubernetes
cluster. It is deployed automatically by the
[Nebari Operator](https://github.com/nebari-dev/nebari-operator) as part of NIC's foundational
software, runs on a schedule as a `CronJob`, and ships each timestamped JSON report to the
built-in web dashboard, a shared PVC, or a ConfigMap.

It exists because answering *"what is actually running on this cluster, where did it come from,
and is it signed?"* should not require manual auditing.

![Provenance Collector dashboard — image inventory with signature, SLSA, SBOM, and update status](https://raw.githubusercontent.com/nebari-dev/nebari-provenance-collector-pack/main/docs/screenshots/dashboard-overview.png)

## What it does

| Capability | Description |
| --- | --- |
| **Image Discovery** | Scans all pods across namespaces, deduplicates by workload owner |
| **Digest Resolution** | Resolves every image tag to its immutable SHA256 digest |
| **Signature Verification** | Checks for cosign signatures (existence or key-based verification) |
| **SLSA Provenance** | Detects SLSA provenance attestations via OCI referrers API |
| **SBOM Detection** | Detects attached SPDX / CycloneDX attestations |
| **Update Checking** | Compares running tags against latest semver tags (configurable level, pre-release filtering) |
| **Helm Release Tracking** | Discovers all deployed Helm releases with chart versions |
| **Web Dashboard** | Optional UI with filters, sorting, pagination, and image detail panel |
| **Grafana Integration** | JSON API compatible with the Infinity datasource for dashboards and alerting |
| **Provenance Reports** | Timestamped JSON reports via the dashboard's upload endpoint (default), a shared PVC, or a ConfigMap, with automatic retention |

{{< callout type="note" >}}
Under active development as part of Nebari Infrastructure Core (NIC). APIs, chart values, and
report schema may change without notice while pre-1.0.
{{< /callout >}}

## Get started

Head to [**Install**](install/) for the operator-managed (recommended) and standalone install
paths.
