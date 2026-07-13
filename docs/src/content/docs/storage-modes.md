---
title: Storage Modes
description: How persistence.mode controls where provenance reports are written and which mode to pick.
---

`persistence.mode` controls how reports get from the collector Job to the
dashboard. Pick one:

| Mode | What it does | Use it when |
|---|---|---|
| `http` *(default)* | Collector POSTs the JSON to the dashboard's internal Service (`{fullname}-web-internal:8081/internal/reports`, cluster-DNS only). The dashboard pod owns a single PVC; nothing else mounts it. | You're on RWO-only storage (Hetzner `csi.hetzner.cloud`, most cloud CSIs). This is the safe default. |
| `pvc` | Collector and dashboard share one PVC at `config.reportPath`. `persistence.storageClass` is required and must be ReadWriteMany-capable (Longhorn, NFS, EFS, …) unless every collector pod is guaranteed to land on the same node as the dashboard. The chart fails to render if `storageClass` is empty in this mode. | You already have an RWX-capable storage class and prefer filesystem semantics. |
| `configmap` | Collector writes the report to a ConfigMap. The dashboard is not used. | Headless / inspection-only installs. |

## Internal upload endpoint

In `http` mode the internal upload endpoint is exposed through a dedicated
Service (`{fullname}-web-internal`) that the NebariApp / public Ingress never
references. Apply a NetworkPolicy restricting it to the collector
ServiceAccount if your cluster supports it.

## Related pages

- [Architecture](/architecture/) — where the storage sink sits in the overall flow.
- [Configuration](/configuration/) — the `PROVENANCE_REPORT_OUTPUT` variable and related chart values.
