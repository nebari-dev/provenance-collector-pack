---
title: Web Dashboard
description: The React dashboard UI, its JSON API, and how to surface provenance data in Grafana.
---

The UI is a standalone **React + TypeScript SPA** (Vite + Tailwind + the Nebari
design system), built into its own nginx image and deployed as a separate
`Deployment`/`Service`. The Go dashboard is **API-only**: nginx serves the SPA
and reverse-proxies `/api/*` to the dashboard over cluster DNS. Enable both:

```yaml
webUI:
  enabled: true       # dashboard API + report-upload endpoint (required in http mode)
frontend:
  enabled: true       # standalone React UI (nginx)
  keycloak:
    url: https://keycloak.<your-domain>   # required: the browser keycloak-js login endpoint
```

The dashboard provides:

- Summary stat cards (`N / M` ratios for Signed, Verified, SLSA, SBOM; absolute counts for Images, Updates, Helm)
- Report timeline to browse historical reports, with an opt-in `+N / -N` unique-image delta badge between adjacent scans (`webUI.features.timelineDeltas`)
- Filterable, sortable, paginated image table (truncated workload column with full name on hover)
- Click any image row for a detail drawer showing signature, SLSA, SBOM, and update info
- Helm releases table
- Light / Dark / System theme, chosen from the profile menu (defaults to System)
- **Run Scan** button — admin-gated; triggers a one-shot Job from the same CronJob template the schedule uses. Hidden unless `webUI.oidcIssuer` is set and the calling user's OIDC groups intersect with `webUI.adminGroups`. Auto-cleanup after `webUI.manualJobTTL` (default 1h).
- **Export** button (`CSV` / `Markdown` / `JSON`) — downloads whichever report is currently selected on the timeline, not just the latest

## Authentication

The SPA runs the OIDC login in the browser via `keycloak-js` (PKCE), attaching
the access token to every `/api` call; nginx forwards it to the dashboard,
which validates it against Keycloak. Under `nebariapp.enabled: true` the
operator provisions the public SPA client and registers routing/landing-page —
the gateway itself does **not** enforce auth
(`nebariapp.auth.enforceAtGateway: false`). The operator also wires
`webUI.oidcIssuer` / `webUI.adminGroups` from the `nebariapp.auth` block so
Run Scan lights up for users in the configured groups. See the
[NebariApp CRD reference](/nebariapp-crd-reference/) for the full field list.

## Branding

The UI ships with built-in Nebari branding (title, logos, favicon, theme
colors) and needs no configuration. Operators can rebrand it **without
rebuilding the image**: branding is delivered at runtime through the same
`/config.json` the SPA already fetches for Keycloak settings, and applied before
React mounts (title, favicon, and theme CSS variables) and in the header (logo).

### Configurable fields

| Field | Description |
|---|---|
| `title` | Browser-tab title. |
| `logoUrl` | Header logo (light mode / default). Absolute `http(s)` URL or root-relative path. |
| `logoUrlDark` | Dark-mode header logo. Falls back to `logoUrl`, then the built-in dark logo. |
| `faviconUrl` | Favicon URL. |
| `theme.light` / `theme.dark` | CSS variable overrides per mode. Supported tokens: `primary`, `primaryForeground`, `background`, `foreground`, `secondary`, `secondaryForeground`, `muted`, `mutedForeground`, `accent`, `accentForeground`, `border`, `ring`, `radius`. |

Every field is optional. Any field left empty uses the built-in Nebari default,
so an unbranded install looks exactly as it does today.

### Kubernetes / Helm

Set `frontend.branding` (and optionally `frontend.title`) in values. The chart
renders them into the `/config.json` ConfigMap mounted into the nginx pod:

```yaml
frontend:
  enabled: true
  title: "Acme Provenance"
  branding:
    logoUrl: "https://cdn.acme.example/logo.svg"
    logoUrlDark: "https://cdn.acme.example/logo-dark.svg"
    faviconUrl: "https://cdn.acme.example/favicon.svg"
    theme:
      light:
        primary: "oklch(55% 0.19 250)"
        primaryForeground: "#ffffff"
      dark:
        primary: "oklch(62% 0.21 250)"
```

A branding-only `helm upgrade` rolls the frontend pod automatically (the
deployment is annotated with a checksum of the rendered ConfigMap).

### Outside Kubernetes

Running the standalone `frontend` image (or the Vite dev server) without a chart,
branding resolves from, in order:

1. A **local `config.json`** — the copy baked into the image, or a file mounted
   over `/usr/share/nginx/html/config.json`, or one pointed to by
   `BRANDING_CONFIG_FILE`.
2. **Environment variables**, overlaid onto that file at container start by the
   image entrypoint (requires the standalone image; a no-op under the read-only
   Kubernetes mount):

   | Env var | Field |
   |---|---|
   | `BRANDING_TITLE` | `title` |
   | `BRANDING_LOGO_URL` | `logoUrl` |
   | `BRANDING_LOGO_URL_DARK` | `logoUrlDark` |
   | `BRANDING_FAVICON_URL` | `faviconUrl` |
   | `BRANDING_THEME` | `theme` (raw JSON, e.g. `'{"light":{"primary":"#0066cc"},"dark":{}}'`) |
   | `KEYCLOAK_URL` / `KEYCLOAK_REALM` / `KEYCLOAK_CLIENT_ID` | `keycloak.*` |

   ```bash
   docker run -p 8080:8080 \
     -e KEYCLOAK_URL=https://kc.acme.example \
     -e BRANDING_TITLE="Acme Provenance" \
     -e BRANDING_LOGO_URL=https://cdn.acme.example/logo.svg \
     ghcr.io/nebari-dev/provenance-collector-pack/frontend
   ```
3. **Built-in Nebari defaults** for any field still unset.

Precedence overall is therefore: chart-rendered `config.json` (in Kubernetes) →
local `config.json` file → `BRANDING_*` env vars → built-in defaults.

### Security

Theme token values are validated in the browser before they are applied: any
value containing CSS-injection characters (`;`, `{`, `}`, `<`, `>`, quotes,
backslash, `url(`, `expression(`, `javascript:`) is dropped rather than injected
into the stylesheet. Logo and favicon URLs are restricted to `http(s)` URLs and
root-relative paths.

## Running the UI locally

Port-forward the dashboard API and point the Vite dev server at it (auth
bypassed for local dev):

```bash
kubectl port-forward svc/provenance-collector-web 8080:8080 -n provenance-system &
cd frontend
npm ci
VITE_DEV_NO_AUTH=true WEBAPI_URL=http://localhost:8080 npm run dev   # → http://localhost:5173
```

The [`dev/Makefile`](https://github.com/nebari-dev/provenance-collector-pack/tree/main/dev)
wraps this as `make ui-up` (install the chart, API only) + `make seed`
(sample reports) + `make ui-dev` (port-forward + Vite).

## Dashboard API

The dashboard exposes a JSON API that the SPA and external tools consume:

| Endpoint | Description |
|---|---|
| `GET /api/reports` | List all reports (newest first) with summary |
| `GET /api/reports/latest` | Get the most recent report |
| `GET /api/reports/<filename>` | Get a specific report by filename |
| `GET /api/export?format=csv\|markdown\|md` | Render the selected report as CSV or Markdown. Optional `&filename=<file>` to pin a historical report; defaults to latest. |
| `GET /api/me` | Calling user's identity + feature flags. Returns `authEnabled`, `canRunScan`, `features.timelineDeltas`. |
| `POST /api/scan` | Trigger a manual scan Job. 403 if the caller isn't in an admin group, 503 if `PROVENANCE_NAMESPACE` / `PROVENANCE_CRONJOB_NAME` aren't configured. |
| `GET /healthz` | Health check |

The JSON payloads follow the structure documented in the
[Report Schema reference](/report-schema/).

## Grafana integration

The provenance data can be surfaced in Grafana using the
[Infinity datasource](https://grafana.com/grafana/plugins/yesoreyeram-infinity-datasource/)
plugin, which queries the dashboard's JSON API.

### Setup

1. Enable the web dashboard (`webUI.enabled: true`)
2. Install the Infinity datasource in Grafana
3. Add a datasource pointing at the dashboard service:
   - **URL:** `http://provenance-collector-web.provenance-system.svc:8080`
   - **Type:** JSON

### Example panels

**Stat panel** (unique images count):

- Type: JSON, URL: `/api/reports/latest`
- Column: `summary.uniqueImages`

**Images table**:

- Type: JSON, URL: `/api/reports/latest`, Root: `images`
- Columns: `image`, `namespace`, `signature.signed`, `provenance.hasProvenance`, `update.updateAvailable`

**Alerting** (images with updates):

```
WHEN count() OF images WHERE updateAvailable = true IS ABOVE 0
```

An example dashboard (11 panels covering unique images, signature status, SLSA
provenance, Helm releases, and more) is available at
[`examples/grafana-dashboard.json`](https://github.com/nebari-dev/provenance-collector-pack/blob/main/examples/grafana-dashboard.json).
Import it directly into Grafana as a `dashboard.grafana.app/v2beta1` resource.
