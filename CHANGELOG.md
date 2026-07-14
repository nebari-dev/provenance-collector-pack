# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-07-13

First stable release. Supersedes the `0.1.0-alpha.*` pre-releases.

### Added
- Core provenance collector: image discovery, digest resolution, cosign
  signature verification (keyless and key-based), SBOM detection, SLSA
  provenance detection, semver update checking, and Helm release tracking.
- Report output modes selected by `persistence.mode`: HTTP upload to the
  dashboard's internal endpoint (default, RWO-safe), a shared PVC, or a
  ConfigMap.
- Web dashboard: a standalone React + TypeScript SPA (served by nginx) backed
  by an API-only Go service, with in-browser OIDC login (`keycloak-js`, PKCE).
  - Summary stat cards, a report timeline with opt-in unique-image delta
    badges, and a filterable/sortable/paginated image table with a detail
    drawer.
  - Report export as CSV, Markdown, or JSON for the selected report.
  - Admin-gated "Run Scan" button that triggers a one-shot Job from the
    CronJob template, with automatic cleanup of manual Jobs.
- SBOM (SPDX) and SLSA provenance (`mode=max`) attestations attached to the
  published collector image on every build, discoverable via the OCI referrers
  API.
- Helm chart: CronJob, RBAC, report storage, dashboard and frontend
  Deployments/Services, and optional NebariApp CRD integration.
- Grafana dashboard example wired to the JSON API via the Infinity datasource.
- Documentation site built with Astro + Starlight and the shared
  `@nebari/starlight` theme, deployed to GitHub Pages.
- SecurityContext hardening (runAsNonRoot, readOnlyRootFilesystem, drop ALL
  capabilities).

### Changed
- README restructured operator-first, with refreshed dashboard sections.
- Configuration reference is generated from a single source of truth
  (`internal/configspec`), guarded against drift in CI.
- Integration test runs on `action-nebari-sandbox` (platform profile, v2)
  instead of a bare kind cluster.
- CI actions bumped to Node-24-compatible majors; releases stamp
  `examples/*.yaml` to the released version.

### Fixed
- SBOM and SLSA provenance detection now read both the OCI referrers index and
  BuildKit's in-index attestation manifests, so attestations attached by
  `docker/build-push-action` are discovered and shown in the dashboard. The
  legacy cosign attestation tag (`.att`) is retained as a fallback for images
  attested with older `cosign attest` runs.

### Known limitations
- The published collector image ships with SBOM and SLSA provenance
  attestations but is **not yet cosign-signed**, so when the collector scans
  its own cluster its image shows as unsigned. Image signing is tracked in #29,
  deferred pending a decision on the signing trust root (keyless Sigstore vs.
  key-managed).
- Air-gapped clusters and private registry mirrors are not yet supported
  (tracked in #1).
