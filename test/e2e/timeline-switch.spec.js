// Timeline scan-switch end-to-end test.
//
// Regression coverage for https://github.com/nebari-dev/nebari-provenance-collector-pack/issues/19:
// clicking a non-latest timeline entry must swap the Images and Helm tables, not just the stat cards.
//
// Preconditions (set up by the calling workflow):
//   - The dashboard is reachable at $DASHBOARD_URL (default http://localhost:8080).
//   - Two provenance reports exist on disk. The newer one must include an
//     image whose name contains $TIMELINE_NEW_IMAGE (default "busybox") that
//     is absent from the older one.

const { test, expect } = require('@playwright/test');

const BASE = process.env.DASHBOARD_URL || 'http://localhost:8080';
const NEW_IMAGE_MARKER = process.env.TIMELINE_NEW_IMAGE || 'busybox';

test('timeline scan-switch updates stat cards AND images table', async ({ page }) => {
  await page.goto(BASE, { waitUntil: 'networkidle' });
  await page.waitForSelector('.timeline-item');

  const items = page.locator('.timeline-item');
  await expect(items, 'expected exactly two reports in the timeline').toHaveCount(2, { timeout: 10_000 });

  // Reports are sorted DESC by generatedAt — item 0 is newest, item 1 is older.
  const newest = items.nth(0);
  const older = items.nth(1);

  // --- Newest report (default state) ----------------------------------------
  await newest.click();
  await page.waitForLoadState('networkidle');

  const newestCount = await readImagesStatCount(page);
  const newestRows = await page.locator('#images-table tbody tr').count();
  const newestHasNewImage = await rowsContaining(page, NEW_IMAGE_MARKER);

  expect(newestRows, 'newest report should have at least one image row').toBeGreaterThan(0);
  expect(newestHasNewImage, `newest report must contain "${NEW_IMAGE_MARKER}"`).toBeGreaterThan(0);

  // --- Older report ---------------------------------------------------------
  await older.click();
  // The click triggers a fetch + re-render; wait for the active class to flip
  // before reading state.
  await expect(older).toHaveClass(/\bactive\b/, { timeout: 5_000 });
  await page.waitForLoadState('networkidle');

  const olderCount = await readImagesStatCount(page);
  const olderRows = await page.locator('#images-table tbody tr').count();
  const olderHasNewImage = await rowsContaining(page, NEW_IMAGE_MARKER);

  // 1. The image stat card must have moved.
  expect(olderCount, 'older report should have a smaller unique-image count than newest')
    .toBeLessThan(newestCount);

  // 2. The Images TABLE must have moved with the stat card.
  // This is what #19 currently breaks: stat updates, table doesn't.
  expect(olderRows, 'older report should have fewer table rows than newest')
    .toBeLessThan(newestRows);

  // 3. The image added between the two scans must not be in the older table.
  expect(olderHasNewImage, `older report must NOT contain "${NEW_IMAGE_MARKER}"`).toBe(0);
});

// readImagesStatCount returns the integer shown in the "Images" stat card.
async function readImagesStatCount(page) {
  const text = await page.locator('#stat-all .value').textContent();
  const n = parseInt((text || '').trim(), 10);
  if (Number.isNaN(n)) {
    throw new Error(`could not parse stat card value: "${text}"`);
  }
  return n;
}

// rowsContaining returns the count of <tbody tr> rows whose text contains needle.
async function rowsContaining(page, needle) {
  return page.locator('#images-table tbody tr').filter({ hasText: needle }).count();
}
