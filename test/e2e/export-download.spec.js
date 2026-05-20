// Export menu download-validation test.
//
// Drives the dashboard's Export dropdown for each format (CSV, Markdown, JSON),
// captures the actual download event, and shape-validates the file against the
// dashboard's own JSON API (the source of truth — the export endpoints render
// from the same on-disk report).
//
// Preconditions (set up by the calling workflow):
//   - The dashboard is reachable at $DASHBOARD_URL (default http://localhost:8080).
//   - At least one provenance report exists on disk.

const { test, expect } = require('@playwright/test');
const fs = require('fs');

const BASE = process.env.DASHBOARD_URL || 'http://localhost:8080';

test.describe('Export menu', () => {
  let truth;
  let truthFilename;

  test.beforeAll(async ({ request }) => {
    // Source of truth: the API response for the latest report. The Export
    // endpoints render from the same file, so any disagreement is a bug in
    // the renderer.
    const list = await request.get(`${BASE}/api/reports`);
    expect(list.ok(), 'GET /api/reports must succeed').toBeTruthy();
    const reports = await list.json();
    expect(reports.length, 'expected at least one report').toBeGreaterThan(0);
    truthFilename = reports[0].filename;
    const latest = await request.get(`${BASE}/api/reports/${truthFilename}`);
    expect(latest.ok(), 'GET /api/reports/<latest> must succeed').toBeTruthy();
    truth = await latest.json();
  });

  test('Export menu toggles open and closed', async ({ page }) => {
    await page.goto(BASE, { waitUntil: 'networkidle' });

    const toggle = page.locator('#export-toggle');
    const menu = page.locator('#export-popover');

    await expect(toggle).toHaveAttribute('aria-expanded', 'false');
    await expect(menu).toBeHidden();

    await toggle.click();
    await expect(toggle).toHaveAttribute('aria-expanded', 'true');
    await expect(menu).toBeVisible();

    // Escape closes.
    await page.keyboard.press('Escape');
    await expect(toggle).toHaveAttribute('aria-expanded', 'false');
    await expect(menu).toBeHidden();

    // Re-open, then click outside to close.
    await toggle.click();
    await expect(menu).toBeVisible();
    await page.locator('#stats').click();
    await expect(menu).toBeHidden();
  });

  test('CSV download has the expected header and one row per image', async ({ page }) => {
    const csv = await openMenuAndDownload(page, '#export-csv');

    const lines = csv.split('\n').filter((l) => l.length > 0);
    const header = lines[0];
    const dataRows = lines.slice(1);

    expect(header, 'CSV header must lead with the expected schema columns')
      .toBe('Image,Namespace,Workload Kind,Workload Name,Digest,Signed,Verified,SLSA Provenance,SBOM,SBOM Format,Update Available,Current Tag,Latest In Major');

    expect(dataRows.length, 'one CSV row per image in the source report')
      .toBe(truth.images.length);

    // Every truth image should appear on exactly one row. Match on the image
    // field (first column) — that's the only one guaranteed unique enough to
    // not collide on a kind cluster with multiple replicasets.
    for (const img of truth.images) {
      const found = dataRows.some((row) => row.startsWith(csvField(img.image) + ','));
      expect(found, `CSV must contain a row for ${img.image}`).toBeTruthy();
    }
  });

  test('Markdown download has the expected sections', async ({ page }) => {
    const md = await openMenuAndDownload(page, '#export-md');

    expect(md).toContain('# Provenance Report');
    expect(md).toContain('## Summary');
    expect(md).toContain('## Container Images');

    // One row per image under "## Container Images". Find the table that
    // follows the heading and count the data rows (skip header + separator).
    const imgSection = md.split('## Container Images')[1];
    expect(imgSection, 'Container Images section must exist').toBeTruthy();
    const imgTable = imgSection.split('## ')[0]; // stop at the next H2 if any
    const imgDataRows = imgTable
      .split('\n')
      .filter((l) => l.startsWith('|') && !l.startsWith('| ---') && !l.includes('| Image |'));
    expect(imgDataRows.length, 'one Markdown image row per image in the source report')
      .toBe(truth.images.length);

    if (truth.helmReleases && truth.helmReleases.length > 0) {
      expect(md).toContain('## Helm Releases');
      const helmSection = md.split('## Helm Releases')[1];
      const helmRows = helmSection
        .split('\n')
        .filter((l) => l.startsWith('|') && !l.startsWith('| ---') && !l.includes('| Release |'));
      expect(helmRows.length, 'one Markdown helm row per release in the source report')
        .toBe(truth.helmReleases.length);
    }
  });

  test('JSON download matches the API response byte-shape', async ({ page }) => {
    const raw = await openMenuAndDownload(page, '#export-json');
    const parsed = JSON.parse(raw);

    // The JSON download is just a passthrough of /api/reports/<filename>, so
    // top-level shape and image / helm counts must match the truth exactly.
    expect(parsed.metadata.generatedAt).toBe(truth.metadata.generatedAt);
    expect(parsed.summary.uniqueImages).toBe(truth.summary.uniqueImages);
    expect(parsed.summary.totalImages).toBe(truth.summary.totalImages);
    expect(parsed.images.length).toBe(truth.images.length);
    expect((parsed.helmReleases || []).length).toBe((truth.helmReleases || []).length);
  });
});

// openMenuAndDownload opens the Export menu, clicks the given menu-item
// selector, captures the download event, and returns the file contents as a
// UTF-8 string.
async function openMenuAndDownload(page, itemSelector) {
  await page.goto(BASE, { waitUntil: 'networkidle' });
  await page.locator('#export-toggle').click();
  await expect(page.locator('#export-popover')).toBeVisible();

  const [download] = await Promise.all([
    page.waitForEvent('download'),
    page.locator(itemSelector).click(),
  ]);
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
