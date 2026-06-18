# Provenance Collector docs site

Hugo site for the pack's user-facing documentation. Consumes
[`nebari-hugo-theme`](https://github.com/nebari-dev/nebari-hugo-theme) as a Hugo Module — this
directory only contains content + a 70-line `hugo.toml`. The theme owns the chrome (header,
sidebar, search, dark-mode toggle, code highlighting, callouts, breadcrumbs, ToC, edit-link).

Built and deployed to GitHub Pages by [`.github/workflows/docs-deploy.yaml`](../.github/workflows/docs-deploy.yaml)
on every push to `main` (and built without publishing on PRs that touch `site/`).

## Local development

```bash
cd site
hugo mod get -u   # first-time: pull the theme module
hugo server       # http://localhost:1313/nebari-provenance-collector-pack/
hugo              # build to public/
```

The `baseURL` in `hugo.toml` is set to the deployed Pages URL, so `hugo server` serves under
that path prefix. Open `http://localhost:1313/nebari-provenance-collector-pack/` (not just `/`).

## Layout

```
site/
  hugo.toml           Theme config — tabs, sidebar tree, editBase, params
  go.mod              Hugo Module init (theme is fetched on `hugo mod get`)
  content/
    _index.md         Overview landing
    install.md        Install runbook (operator-managed + standalone)
```

## Adding a page

Drop a `.md` under `content/`. Add a `[[params.sidebar.items]]` entry in `hugo.toml` to surface
it in the left sidebar. The theme's ToC widget on the right is auto-generated from the H2/H3
headings in the page.

## What the theme provides

See the [upstream README](https://github.com/nebari-dev/nebari-hugo-theme#whats-shipped) for
the full feature list — search, code-block copy buttons, Mermaid diagrams, breadcrumbs,
last-updated stamps, i18n, versioning, callout shortcodes, themed 404, responsive layout, etc.

## Why a separate Hugo site?

The reference markdowns under `../docs/` (`configuration.md`, `nebariapp-crd-reference.md`,
`report-schema.md`) are CI-generated artifacts that the main `README.md` links to. The Hugo
site under `site/` is the *human-readable* docs surface. If/when this skeleton graduates into
[`nebari-software-pack-template`](https://github.com/nebari-dev/nebari-software-pack-template),
every Nebari software pack inherits the same shape: `site/` for narrative docs, `docs/` for
auto-generated references.
