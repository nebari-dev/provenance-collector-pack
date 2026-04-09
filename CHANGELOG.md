# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Core provenance collector: image discovery, digest resolution, signature
  verification, SBOM detection, update checking, and Helm release tracking
- PVC and ConfigMap report output writers
- Key-based cosign signature verification using `cosign.VerifyImageSignatures`
- SBOM format detection via in-toto attestation JSON parsing
- Optional web dashboard with report timeline, summary stats, image and Helm
  release tables
- Helm chart with CronJob, RBAC, PVC, optional web dashboard Deployment/Service,
  and NebariApp CRD integration
- CI workflows: lint, test, integration test (kind), image build (GHCR), and
  release (Helm chart to nebari-dev/helm-repository)
- SecurityContext hardening (runAsNonRoot, readOnlyRootFilesystem, drop ALL
  capabilities)
