+++
title = 'Install'
description = "Two paths to install the Provenance Collector pack: operator-managed default, and standalone for vanilla Kubernetes."
+++

The Provenance Collector is normally installed by the Nebari Operator as part of NIC's
foundational software — you don't run `helm` yourself, the operator and ArgoCD do it for you.

## Operator-managed install (default)

A complete ArgoCD `Application` manifest lives at
[`examples/argocd-application.yaml`](https://github.com/nebari-dev/nebari-provenance-collector-pack/blob/main/examples/argocd-application.yaml)
and is auto-stamped to the latest released chart version on every release. Apply it from your
gitops repo or directly:

```bash
kubectl apply -f examples/argocd-application.yaml
```

The values most users adjust:

```yaml
nebariapp:
  enabled: true                       # register the pack with the Nebari Operator
  hostname: provenance.<your-domain>  # public URL the dashboard responds on

webUI:
  enabled: true                       # default true; required when persistence.mode=http
  features:
    timelineDeltas: false             # opt-in; show +N/-N badges between scans
```

### Verify

```bash
kubectl get application provenance-collector -n argocd
kubectl get cronjob -n provenance-system
kubectl get pods -n provenance-system -l app.kubernetes.io/name=provenance-collector
```

The operator wires routing (Envoy Gateway), OIDC (Keycloak), and surfaces the dashboard on the
Nebari Landing page — no extra config needed beyond `hostname`.

## Standalone install (without the Nebari Operator)

{{< callout type="warning" >}}
Use this path only on a vanilla Kubernetes cluster *without* NIC. Without the operator you're
responsible for routing and OIDC yourself if you want the dashboard reachable from outside the
cluster.
{{< /callout >}}

### Prerequisites

| Tool | Minimum version | Notes |
| --- | --- | --- |
| `kubectl` | 1.26+ | Cluster interaction |
| `helm` | 3.14+ | Chart install |
| Kubernetes cluster | 1.26+ | Local (kind / k3d / minikube) or remote |
| Cluster permissions | `cluster-admin` | Chart creates a `ClusterRole` + `ClusterRoleBinding` |

### Install

```bash
helm repo add nebari https://nebari-dev.github.io/helm-repository
helm repo update

helm install provenance-collector nebari/provenance-collector \
  --namespace provenance-system \
  --create-namespace
```

See the
[repo README](https://github.com/nebari-dev/nebari-provenance-collector-pack#standalone-install-without-the-nebari-operator)
for the full standalone runbook (verify, trigger a manual run, view the report, uninstall).
