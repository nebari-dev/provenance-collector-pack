---
title: Verifying the Collector Image
description: Verify the cosign signature, SLSA provenance, and SBOM attached to the published provenance-collector images.
---

The published collector and dashboard-frontend images are signed with **keyless
cosign** (Sigstore) using the release workflow's GitHub Actions OIDC identity -
there are no long-lived signing keys. Every build also attaches an SPDX SBOM and
a SLSA provenance attestation as OCI referrers.

## Verify the signature

Requires [cosign](https://docs.sigstore.dev/) v3+. **Identity pinning is
mandatory**: without `--certificate-identity-regexp` and
`--certificate-oidc-issuer`, cosign would accept a signature from any identity,
which defeats the purpose.

```bash
cosign verify \
  --certificate-identity-regexp '^https://github.com/nebari-dev/provenance-collector-pack/\.github/workflows/build-image\.yaml@refs/tags/v.*$' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  quay.io/nebari/provenance-collector:0.1.0
```

The same digest is published to GHCR - swap the reference to verify it, or the
dashboard UI image:

- `ghcr.io/nebari-dev/provenance-collector-pack:0.1.0`
- `ghcr.io/nebari-dev/provenance-collector-pack/frontend:0.1.0`

A successful run prints the verified certificate's identity and OIDC issuer.

## Inspect the SBOM and provenance

```bash
cosign tree quay.io/nebari/provenance-collector:0.1.0
```

This lists the attached SPDX SBOM and SLSA provenance attestations. The
collector detects these same referrers when it scans a cluster, which is why
its own image shows **Signed / SBOM / SLSA** in the dashboard rather than the
grey "no provenance" state.
