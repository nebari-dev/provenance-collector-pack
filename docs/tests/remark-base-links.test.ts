import { describe, it, expect } from 'vitest';
import { remark } from 'remark';
import remarkBaseLinks, { prefixUrl } from '../src/plugins/remark-base-links';

describe('prefixUrl', () => {
  const cases: Array<{ name: string; url: string; base: string; want: string }> = [
    { name: 'base "/" leaves links unchanged', url: '/configuration/', base: '/', want: '/configuration/' },
    { name: 'sub-path base prefixes internal links', url: '/configuration/', base: '/provenance-collector-pack/', want: '/provenance-collector-pack/configuration/' },
    { name: 'prefixes image paths', url: '/screenshots/dashboard-overview.png', base: '/provenance-collector-pack/', want: '/provenance-collector-pack/screenshots/dashboard-overview.png' },
    { name: 'never rewrites external links', url: 'https://nebari.dev', base: '/provenance-collector-pack/', want: 'https://nebari.dev' },
    { name: 'never rewrites protocol-relative links', url: '//example.com/x', base: '/provenance-collector-pack/', want: '//example.com/x' },
    { name: 'never rewrites anchor-only links', url: '#section', base: '/provenance-collector-pack/', want: '#section' },
    { name: 'preserves anchors on internal links', url: '/configuration/#collector', base: '/provenance-collector-pack/', want: '/provenance-collector-pack/configuration/#collector' },
    { name: 'idempotent on already-prefixed links', url: '/provenance-collector-pack/configuration/', base: '/provenance-collector-pack/', want: '/provenance-collector-pack/configuration/' },
  ];
  for (const c of cases) {
    it(c.name, () => {
      expect(prefixUrl(c.url, c.base)).toBe(c.want);
    });
  }
});

describe('remarkBaseLinks plugin', () => {
  it('rewrites link and image urls in a markdown document', async () => {
    const md = 'See [Configuration](/configuration/) and ![img](/img/a.png) and [ext](https://nebari.dev)';
    const out = String(
      await remark().use(remarkBaseLinks, { base: '/provenance-collector-pack/' }).process(md),
    );
    expect(out).toContain('(/provenance-collector-pack/configuration/)');
    expect(out).toContain('(/provenance-collector-pack/img/a.png)');
    expect(out).toContain('(https://nebari.dev)');
  });

  it('is a no-op when base is "/"', async () => {
    const md = '[C](/configuration/)';
    const out = String(await remark().use(remarkBaseLinks, { base: '/' }).process(md));
    expect(out).toContain('(/configuration/)');
  });
});
