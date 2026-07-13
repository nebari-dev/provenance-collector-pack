// Export menu download-validation test.
//
// Drives the dashboard's Export dropdown for each format (CSV, Markdown, JSON),
// captures the actual download event, and shape-validates the file against the
// dashboard's own JSON API (the source of truth — the export endpoints render
// from the same on-disk report).
//
// Preconditions (set up by the calling workflow):
//   - The React SPA is reachable at $DASHBOARD_URL (Vite, proxying /api to the
//     dashboard). Auth is injected via window.__PW_E2E_AUTH__.
//   - At least one provenance report exists on disk.

const { test, expect } = require('@playwright/test');
const fs = require('fs');

const BASE = process.env.DASHBOARD_URL || 'http://localhost:5173';

async function injectAuth(page) {
  await page.addInitScript(() => {
    window.__PW_E2E_AUTH__ = {
      authenticated: true,
      token: 'e2e-token',
      idTokenParsed: { name: 'E2E', email: 'e2e@example.com', preferred_username: 'e2e' },
    };
  });
}

test.describe('Export menu', () => {
  let truth;
  let truthFilename;

  test.beforeAll(async ({ request }) => {
    // Source of truth: the API response for the latest report. The Export
    // endpoints render from the same file, so any disagreement is a bug in the
    // renderer. The API request goes through Vite's /api proxy to the dashboard.
    const list = await request.get(`${BASE}/api/reports`);
    expect(list.ok(), 'GET /api/reports must succeed').toBeTruthy();
    const reports = await list.json();
    expect(reports.length, 'expected at least one report').toBeGreaterThan(0);
    truthFilename = reports[0].filename;
    const latest = await request.get(`${BASE}/api/reports/${truthFilename}`);
    expect(latest.ok(), 'GET /api/reports/<latest> must succeed').toBeTruthy();
    truth = await latest.json();
  });

  test.beforeEach(async ({ page }) => {
    await injectAuth(page);
  });

  test('Export menu toggles open and closed', async ({ page }) => {
    await page.goto(BASE, { waitUntil: 'networkidle' });

    const trigger = page.getByRole('button', { name: 'Export' });
    const csvItem = page.getByRole('menuitem', { name: 'CSV' });

    await expect(csvItem).toHaveCount(0);

    await trigger.click();
    await expect(csvItem).toBeVisible();

    // Escape closes.
    await page.keyboard.press('Escape');
    await expect(csvItem).toHaveCount(0);

    // Re-open, then click outside (the page title) to close.
    await trigger.click();
    await expect(csvItem).toBeVisible();
    await page.getByRole('heading', { name: 'Provenance' }).click();
    await expect(csvItem).toHaveCount(0);
  });

  test('CSV download has the expected header and one row per image', async ({ page }) => {
    const csv = await openMenuAndDownload(page, 'CSV');

    const lines = csv.split('\n').filter((l) => l.length > 0);
    const header = lines[0];
    const dataRows = lines.slice(1);

    expect(header, 'CSV header must lead with the expected schema columns').toBe(
      'Image,Namespace,Workload Kind,Workload Name,Digest,Signed,Verified,SLSA Provenance,SBOM,SBOM Format,Update Available,Current Tag,Latest In Major',
    );

    expect(dataRows.length, 'one CSV row per image in the source report').toBe(truth.images.length);

    for (const img of truth.images) {
      const found = dataRows.some((row) => row.startsWith(csvField(img.image) + ','));
      expect(found, `CSV must contain a row for ${img.image}`).toBeTruthy();
    }
  });

  test('Markdown download has the expected sections', async ({ page }) => {
    const md = await openMenuAndDownload(page, 'Markdown');

    expect(md).toContain('# Provenance Report');
    expect(md).toContain('## Summary');
    expect(md).toContain('## Container Images');

    const imgSection = md.split('## Container Images')[1];
    expect(imgSection, 'Container Images section must exist').toBeTruthy();
    const imgTable = imgSection.split('## ')[0];
    const imgDataRows = markdownDataRows(imgTable, '| Image |');
    expect(imgDataRows.length, 'one Markdown image row per image in the source report').toBe(
      truth.images.length,
    );

    if (truth.helmReleases && truth.helmReleases.length > 0) {
      expect(md).toContain('## Helm Releases');
      const helmSection = md.split('## Helm Releases')[1];
      const helmRows = markdownDataRows(helmSection, '| Release |');
      expect(helmRows.length, 'one Markdown helm row per release in the source report').toBe(
        truth.helmReleases.length,
      );
    }
  });

  test('JSON download matches the API response byte-shape', async ({ page }) => {
    const raw = await openMenuAndDownload(page, 'JSON');
    const parsed = JSON.parse(raw);

    expect(parsed.metadata.generatedAt).toBe(truth.metadata.generatedAt);
    expect(parsed.summary.uniqueImages).toBe(truth.summary.uniqueImages);
    expect(parsed.summary.totalImages).toBe(truth.summary.totalImages);
    expect(parsed.images.length).toBe(truth.images.length);
    expect((parsed.helmReleases || []).length).toBe((truth.helmReleases || []).length);
  });
});

// openMenuAndDownload opens the Export menu, clicks the given menu item by its
// accessible name, captures the download event, and returns the file contents.
async function openMenuAndDownload(page, itemName) {
  await page.goto(BASE, { waitUntil: 'networkidle' });
  await page.getByRole('button', { name: 'Export' }).click();
  const item = page.getByRole('menuitem', { name: itemName });
  await expect(item).toBeVisible();

  const [download] = await Promise.all([page.waitForEvent('download'), item.click()]);
  const path = await download.path();
  return fs.readFileSync(path, 'utf-8');
}

// csvField mirrors the server's csvEscape so we can compare predictably. Only
// covers the inputs we actually test against (image names).
function csvField(s) {
  if (/[,"\n]/.test(s)) {
    return '"' + s.replace(/"/g, '""') + '"';
  }
  return s;
}

// markdownDataRows returns the data rows of a Markdown table from a section
// blob, skipping the header (matched by headerNeedle) and the separator row.
function markdownDataRows(section, headerNeedle) {
  return section
    .split('\n')
    .filter((l) => l.startsWith('|'))
    .filter((l) => !l.includes(headerNeedle))
    .filter((l) => !/^\|[-:|\s]+$/.test(l));
}
