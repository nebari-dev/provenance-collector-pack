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
- Published collector and dashboard images are signed with keyless cosign
  (Sigstore, via GitHub Actions OIDC - no managed key) and carry SPDX SBOM and
  SLSA provenance (`mode=max`) attestations, all discoverable via the OCI
  referrers API. See "Verifying the Collector Image" in the docs.
- Helm chart: CronJob, RBAC, report storage, dashboard and frontend
  Deployments/Services, and optional NebariApp CRD integration.
- Grafana dashboard example wired to the JSON API via the Infinity datasource.
- Documentation site built with Astro + Starlight and the shared
  `@nebari/starlight` theme, deployed to Cloudflare Pages and routed through
  `packs.nebari.dev/provenance-collector-pack/`, with per-PR previews.
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
- Air-gapped clusters and private registry mirrors are not yet supported
  (tracked in #1).
