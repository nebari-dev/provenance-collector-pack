// Timeline scan-switch end-to-end test.
//
// Regression coverage for https://github.com/nebari-dev/provenance-collector-pack/issues/19:
// clicking a non-latest timeline entry must swap the Images and Helm tables,
// not just the stat cards.
//
// Preconditions (set up by the calling workflow):
//   - The React SPA is reachable at $DASHBOARD_URL (Vite, proxying /api to the
//     dashboard). Auth is injected via window.__PW_E2E_AUTH__.
//   - Two provenance reports exist on disk. The newer one must include an image
//     whose name contains $TIMELINE_NEW_IMAGE (default "traefik/whoami") that is
//     absent from the older one.

const { test, expect } = require('@playwright/test');

const BASE = process.env.DASHBOARD_URL || 'http://localhost:5173';
const NEW_IMAGE_MARKER = process.env.TIMELINE_NEW_IMAGE || 'traefik/whoami';

async function injectAuth(page) {
  await page.addInitScript(() => {
    window.__PW_E2E_AUTH__ = {
      authenticated: true,
      token: 'e2e-token',
      idTokenParsed: { name: 'E2E', email: 'e2e@example.com', preferred_username: 'e2e' },
    };
  });
}

test('timeline scan-switch updates stat cards AND images table', async ({ page }) => {
  await injectAuth(page);
  await page.goto(BASE, { waitUntil: 'networkidle' });

  const items = page.getByTestId('timeline-item');
  await expect(items, 'expected exactly two reports in the timeline').toHaveCount(2, {
    timeout: 10_000,
  });

  // Reports are sorted DESC by generatedAt — item 0 is newest, item 1 is older.
  const newest = items.nth(0);
  const older = items.nth(1);

  // --- Newest report (default state) ----------------------------------------
  await newest.click();
  await expect(newest).toHaveAttribute('aria-pressed', 'true');
  const newestStat = await readStatAll(page);
  const newestTotal = await readImagesTotal(page);
  const newestHasNewImage = await searchMatchCount(page, NEW_IMAGE_MARKER);

  expect(newestTotal, 'newest report should expose a non-zero image count').toBeGreaterThan(0);
  expect(newestHasNewImage, `newest report must contain "${NEW_IMAGE_MARKER}"`).toBeGreaterThan(0);

  // --- Older report ---------------------------------------------------------
  await older.click();
  await expect(older).toHaveAttribute('aria-pressed', 'true', { timeout: 5_000 });
  const olderStat = await readStatAll(page);
  const olderTotal = await readImagesTotal(page);
  const olderHasNewImage = await searchMatchCount(page, NEW_IMAGE_MARKER);

  // 1. The Images stat card (unique image count) must have moved.
  expect(olderStat, 'older report should have a smaller unique-image stat than newest').toBeLessThan(
    newestStat,
  );

  // 2. The Images TABLE must have moved with the stat card. We compare the
  // section's "(N)" count label (report.images.length) so pagination can't mask
  // the swap — this is what #19 reports as broken: stat updates, table doesn't.
  expect(olderTotal, 'older report should have a smaller table image total than newest').toBeLessThan(
    newestTotal,
  );

  // 3. The image added between the two scans must not be in the older table.
  expect(olderHasNewImage, `older report must NOT contain "${NEW_IMAGE_MARKER}"`).toBe(0);
});

// readStatAll returns the integer in the "Images" stat card (summary.uniqueImages).
async function readStatAll(page) {
  const text = await page.getByTestId('stat-all').getByTestId('stat-value').textContent();
  const n = parseInt((text || '').trim(), 10);
  if (Number.isNaN(n)) throw new Error(`could not parse stat-all value: "${text}"`);
  return n;
}

// readImagesTotal returns the integer in the Container Images section header
// ("(N)" = report.images.length, before pagination).
async function readImagesTotal(page) {
  const text = await page.getByTestId('images-total').textContent();
  const match = (text || '').match(/(\d+)/);
  if (!match) throw new Error(`could not parse images-total: "${text}"`);
  return parseInt(match[1], 10);
}

// searchMatchCount types the needle into the image search box and reads the
// "N of M match" line, returning N. Clears the box afterwards so subsequent
// reads start clean. Avoids driving the Base UI page-size Select.
async function searchMatchCount(page, needle) {
  const search = page.getByTestId('image-search');
  await search.fill(needle);
  const line = page.getByTestId('images-match');
  await expect(line).toBeVisible({ timeout: 5_000 });
  const text = (await line.textContent()) || '';
  await search.fill('');
  const match = text.match(/(\d+)\s+of\s+\d+\s+match/);
  return match ? parseInt(match[1], 10) : 0;
}
