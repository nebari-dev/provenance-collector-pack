# Provenance Collector docs site

Astro 5 + Tailwind v4 + the `@nebari/*` [shadcn registry](https://github.com/nebari-dev/nebari-design). Deployed to GitHub Pages on every push to `main`.

This directory is a **spike** for [`nebari-software-pack-template#14`](https://github.com/nebari-dev/nebari-software-pack-template/issues/14) — the docs-site stack the template should default to. If it lands well, this skeleton graduates into the template's `docs-site/` so every future pack inherits it.

## Local development

```bash
cd docs-site
npm install      # or `bun install` if you have bun
npm run dev      # http://localhost:4321
npm run build    # static output → docs-site/dist
npm run preview  # serve the built output
```

Production build outputs static files under `dist/` — no SSR runtime, no server, no special hosting.

## Layout

```
docs-site/
  astro.config.mjs           Astro + React + Tailwind v4 wiring
  components.json            shadcn config with the @nebari registry
  package.json
  tsconfig.json
  src/
    components/
      ui/                    nebari-design components land here via shadcn add
        button.tsx
        spinner.tsx
    layouts/
      PackDocsLayout.astro
    lib/
      utils.ts               @nebari/utils — cn() helper
    pages/
      index.astro
      install.mdx            MDX page demonstrating @nebari/button use
    styles/
      globals.css            @nebari/theme tokens + Tailwind v4 @theme inline
  README.md
```

## Adding a component from `@nebari`

Browse the registry: <https://nebari-dev.github.io/nebari-design/>. Then run:

```bash
npx shadcn@latest add @nebari/<component-name>
```

The CLI reads `components.json`, fetches the component source from `https://nebari-dev.github.io/nebari-design/r/<component-name>.json`, and lands the `.tsx` file under `src/components/ui/`. The same applies to `@nebari/theme` for the CSS tokens.

Re-running the same command pulls the latest version of the component from upstream — that's how the design system stays in sync with consumers.

## Adding a content page

Static / mostly-prose pages: drop a `.astro` file under `src/pages/`.

Mixed markdown + component pages: drop a `.mdx` file under `src/pages/`. Use `layout: '@/layouts/PackDocsLayout.astro'` in the frontmatter so the page picks up the header / footer chrome. Import components from `@/components/ui/<name>` and render them as JSX inline with the markdown.

```mdx
---
layout: '@/layouts/PackDocsLayout.astro'
title: 'Page title'
---

import { Button } from '@/components/ui/button';

# Heading

Prose...

<Button>I'm a real nebari-design component</Button>
```

## Why this stack

See [`nebari-software-pack-template#14`](https://github.com/nebari-dev/nebari-software-pack-template/issues/14) for the full comparison vs Hugo / Docusaurus and the reasoning behind picking Astro. Short version: Astro is the only static-site generator that consumes nebari-design's shadcn registry **natively** — same `bunx shadcn add @nebari/<name>` flow as [`jbouder/nebari-ui-demo`](https://github.com/jbouder/nebari-ui-demo) or any other React-first consumer in the Nebari ecosystem. Hugo can consume the CSS tokens but not the React components; the divergence would force us to maintain a parallel Hugo theme forever.

## Deployment

The workflow at [`.github/workflows/docs-deploy.yaml`](../.github/workflows/docs-deploy.yaml) builds on every push to `main` (and on PRs that touch `docs-site/`) and publishes to GitHub Pages via the official `actions/deploy-pages` action. PRs build without publishing so reviewers see green CI without touching the live site.

The site lives at <https://nebari-dev.github.io/nebari-provenance-collector-pack/> once Pages is enabled in the repo settings.
