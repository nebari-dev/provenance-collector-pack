## What's New

-

## Installation

```bash
helm repo add nebari https://nebari-dev.github.io/helm-repository/
helm repo update
helm install provenance-collector nebari/provenance-collector \
  --version ${VERSION} \
  --namespace provenance-system \
  --create-namespace \
  --set webUI.enabled=true
```

<details>
<summary>Other install methods</summary>

### From OCI (quay.io)

```bash
helm install provenance-collector \
  oci://quay.io/nebari/charts/provenance-collector \
  --version ${VERSION} \
  --namespace provenance-system \
  --create-namespace \
  --set webUI.enabled=true
```

### From source

```bash
git clone https://github.com/nebari-dev/provenance-collector.git
cd provenance-collector
git checkout v${VERSION}
helm install provenance-collector chart/ \
  --namespace provenance-system \
  --create-namespace \
  --set webUI.enabled=true
```

### Container images

```
quay.io/nebari/provenance-collector:${VERSION}
ghcr.io/nebari-dev/provenance-collector:${VERSION}
```

</details>

> **Note:** This is a pre-release. APIs and report format may change before v1.0.
