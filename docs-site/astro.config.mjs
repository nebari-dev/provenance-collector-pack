// @ts-check
import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';
import react from '@astrojs/react';
import tailwindcss from '@tailwindcss/vite';

// The site is published to GitHub Pages at the repo's default subpath.
// Update `site` + `base` if a custom domain is wired up.
export default defineConfig({
  site: 'https://nebari-dev.github.io',
  base: '/nebari-provenance-collector-pack',
  trailingSlash: 'always',
  integrations: [mdx(), react()],
  // Shiki theme matches Amit's Hugo spike (Catppuccin Mocha via Chroma).
  // Always-dark code regardless of page theme — same palette across the
  // Nebari docs surface.
  markdown: {
    shikiConfig: { theme: 'catppuccin-mocha' },
  },
  vite: {
    plugins: [tailwindcss()],
  },
});
