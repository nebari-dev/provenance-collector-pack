// Timeline scan-switch end-to-end test.
//
// Regression coverage for https://github.com/nebari-dev/nebari-provenance-collector-pack/issues/19:
// clicking a non-latest timeline entry must swap the Images and Helm tables, not just the stat cards.
//
// Preconditions (set up by the calling workflow):
//   - The dashboard is reachable at $DASHBOARD_URL (default http://localhost:8080).
//   - Two provenance reports exist on disk. The newer one must include an
//     image whose name contains $TIMELINE_NEW_IMAGE (default "traefik/whoami") that
//     is absent from the older one.

const { test, expect } = require('@playwright/test');

const BASE = process.env.DASHBOARD_URL || 'http://localhost:8080';
const NEW_IMAGE_MARKER = process.env.TIMELINE_NEW_IMAGE || 'traefik/whoami';

test('timeline scan-switch updates stat cards AND images table', async ({ page }) => {
  await page.goto(BASE, { waitUntil: 'networkidle' });
  await page.waitForSelector('.timeline-item');

  const items = page.locator('.timeline-item');
  await expect(items, 'expected exactly two reports in the timeline').toHaveCount(2, { timeout: 10_000 });

  // The table paginates at 25 rows by default. On a platform-stack cluster the
  // unique image count is ~30, so the sentinel row (alphabetically "t") can
  // land on page 2 — the assertion below would then see 0 occurrences even
  // though the table did swap correctly. Bump page size to 100 so the entire
  // image list fits on one page for both reports.
  await selectMaxPageSize(page);

  // Reports are sorted DESC by generatedAt — item 0 is newest, item 1 is older.
  const newest = items.nth(0);
  const older = items.nth(1);

  // --- Newest report (default state) ----------------------------------------
  await newest.click();
  await page.waitForLoadState('networkidle');
  await selectMaxPageSize(page);

  const newestStat = await readImagesStatCount(page);
  const newestTotal = await readImagesTableTotal(page);
  const newestHasNewImage = await rowsContaining(page, NEW_IMAGE_MARKER);

  expect(newestTotal, 'newest report should expose a non-zero image count').toBeGreaterThan(0);
  expect(newestHasNewImage, `newest report must contain "${NEW_IMAGE_MARKER}"`).toBeGreaterThan(0);

  // --- Older report ---------------------------------------------------------
  await older.click();
  // The click triggers a fetch + re-render; wait for the active class to flip
  // before reading state.
  await expect(older).toHaveClass(/\bactive\b/, { timeout: 5_000 });
  await page.waitForLoadState('networkidle');
  await selectMaxPageSize(page);

  const olderStat = await readImagesStatCount(page);
  const olderTotal = await readImagesTableTotal(page);
  const olderHasNewImage = await rowsContaining(page, NEW_IMAGE_MARKER);

  // 1. The Images stat card (unique image count) must have moved.
  expect(olderStat, 'older report should have a smaller unique-image stat than newest')
    .toBeLessThan(newestStat);

  // 2. The Images TABLE must have moved with the stat card. We compare the
  // panel's "N total" label rather than rendered DOM rows so pagination
  // doesn't mask the swap — the label reflects `r.images.length` directly.
  // This is what #19 reports as broken: stat updates, table doesn't.
  expect(olderTotal, 'older report should have a smaller table image total than newest')
    .toBeLessThan(newestTotal);

  // 3. The image added between the two scans must not be in the older table.
  expect(olderHasNewImage, `older report must NOT contain "${NEW_IMAGE_MARKER}"`).toBe(0);
});

// Timeline-delta badge is an opt-in feature gated on
// /api/me.features.timelineDeltas. The integration cluster's chart install
// leaves it off (matches the chart default), so the live /api/me returns
// `features: { timelineDeltas: false }` and no badges render. We assert both
// states by leaving the live API alone for the "off" case and route-mocking
// /api/me for the "on" case.

test('timeline cards omit deltas when the feature flag is disabled (default)', async ({ page }) => {
  await page.goto(BASE, { waitUntil: 'networkidle' });
  await page.waitForSelector('.timeline-item');

  const items = page.locator('.timeline-item');
  await expect(items).toHaveCount(2, { timeout: 10_000 });

  // No delta badges anywhere — the feature flag defaults to off.
  await expect(page.locator('.timeline-item .delta')).toHaveCount(0);
});

test('timeline cards show unique-image deltas when the feature flag is enabled', async ({ page }) => {
  // Mock /api/me to flip timelineDeltas on. Server-side enforcement is covered
  // by internal/dashboard/server_test.go; this spec exercises the JS gate.
  await page.route('**/api/me', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        authEnabled: false,
        canRunScan: false,
        features: { timelineDeltas: true },
      }),
    });
  });

  await page.goto(BASE, { waitUntil: 'networkidle' });
  await page.waitForSelector('.timeline-item');

  const items = page.locator('.timeline-item');
  await expect(items).toHaveCount(2, { timeout: 10_000 });

  // Newest card has a delta (older neighbor exists) and must be positive,
  // because the sentinel deploy added at least one unique image.
  const newestDelta = items.nth(0).locator('.delta');
  await expect(newestDelta, 'newest card must render a delta').toHaveCount(1);
  await expect(newestDelta).toHaveClass(/\bdelta-up\b/);
  const newestDeltaValue = parseInt(await newestDelta.getAttribute('data-delta') || '0', 10);
  expect(newestDeltaValue, 'newest delta should be a positive image gain').toBeGreaterThan(0);

  // Oldest card has no neighbor to diff against — no delta badge rendered.
  const olderDelta = items.nth(1).locator('.delta');
  await expect(olderDelta, 'oldest card must render no delta').toHaveCount(0);
});

// selectMaxPageSize bumps the table page size to 100 via the pagination
// dropdown if it is rendered. The dropdown only appears when the report has
// more rows than the current page size, so on small reports this is a no-op.
async function selectMaxPageSize(page) {
  const select = page.locator('#images-pagination select');
  if (await select.count()) {
    await select.selectOption('100');
  }
}

// readImagesStatCount returns the integer shown in the "Images" stat card —
// i.e. `summary.uniqueImages` (deduplicated by workload).
async function readImagesStatCount(page) {
  const text = await page.locator('#stat-all .value').textContent();
  const n = parseInt((text || '').trim(), 10);
  if (Number.isNaN(n)) {
    throw new Error(`could not parse stat card value: "${text}"`);
  }
  return n;
}

// readImagesTableTotal returns the integer shown in the Container Images
// panel header (`r.images.length`, i.e. every per-pod entry before pagination).
async function readImagesTableTotal(page) {
  const text = await page.locator('#image-count').textContent();
  const match = (text || '').match(/(\d+)/);
  if (!match) {
    throw new Error(`could not parse #image-count: "${text}"`);
  }
  return parseInt(match[1], 10);
}

// rowsContaining returns the count of <tbody tr> rows whose text contains needle.
async function rowsContaining(page, needle) {
  return page.locator('#images-table tbody tr').filter({ hasText: needle }).count();
}
