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
  vite: {
    plugins: [tailwindcss()],
  },
});
