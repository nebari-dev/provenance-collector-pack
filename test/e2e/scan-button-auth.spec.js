// Run Scan button visibility based on the authenticated user's groups.
//
// The button is rendered hidden by default (display:none in the HTML) and is
// only revealed when `/api/me` returns `canRunScan: true`. This spec exercises
// that frontend gate by mocking the `/api/me` response — no real OIDC issuer
// involved. Server-side enforcement (the bearer → userinfo → admin-group
// check) is covered by internal/dashboard/auth_test.go and
// internal/dashboard/internal_server_test.go.
//
// The chart in the integration workflow is installed with auth disabled, so
// the real `/api/me` would return `{authEnabled: false, canRunScan: false}`.
// The mock here overrides that response *before* navigation so every test
// case sees a deterministic state.

const { test, expect } = require('@playwright/test');

const BASE = process.env.DASHBOARD_URL || 'http://localhost:8080';

// mockMe sets up a route handler that returns the given /api/me payload.
// Must be called BEFORE page.goto() so the response is in place when the
// dashboard's init() fetches it on load.
async function mockMe(page, payload) {
  await page.route('**/api/me', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(payload),
    });
  });
}

test.describe('Run Scan button visibility', () => {
  test('hidden when auth is disabled', async ({ page }) => {
    await mockMe(page, { authEnabled: false, canRunScan: false });
    await page.goto(BASE, { waitUntil: 'networkidle' });

    await expect(page.locator('#btn-scan')).toBeHidden();
  });

  test('hidden when user is authenticated but lacks an admin group', async ({ page }) => {
    await mockMe(page, {
      authEnabled: true,
      email: 'alice@example.com',
      groups: ['platform-users'],
      canRunScan: false,
    });
    await page.goto(BASE, { waitUntil: 'networkidle' });

    await expect(page.locator('#btn-scan')).toBeHidden();
  });

  test('visible when user is in an admin group', async ({ page }) => {
    await mockMe(page, {
      authEnabled: true,
      email: 'admin@example.com',
      groups: ['provenance-admins'],
      canRunScan: true,
    });
    await page.goto(BASE, { waitUntil: 'networkidle' });

    const btn = page.locator('#btn-scan');
    await expect(btn).toBeVisible();
    // Button is not disabled in its idle state; the click handler manages
    // disabled-while-polling on its own.
    await expect(btn).toBeEnabled();
  });

  test('clicking the button posts to /api/scan', async ({ page }) => {
    await mockMe(page, {
      authEnabled: true,
      email: 'admin@example.com',
      groups: ['provenance-admins'],
      canRunScan: true,
    });

    // Stub /api/scan with a successful response so we don't kick off a real
    // Job and so the polling loop in the page can settle quickly.
    let scanPosts = 0;
    await page.route('**/api/scan', async (route) => {
      if (route.request().method() === 'POST') {
        scanPosts++;
        await route.fulfill({
          status: 202,
          contentType: 'application/json',
          body: JSON.stringify({ jobName: 'manual-scan-spec-stub' }),
        });
        return;
      }
      await route.continue();
    });

    await page.goto(BASE, { waitUntil: 'networkidle' });
    await page.locator('#btn-scan').click();

    // The handler awaits the POST before flipping the button label, so by the
    // time the label changes we know the POST went out.
    await expect(page.locator('#btn-scan')).toContainText(/scan running/i, { timeout: 5_000 });
    expect(scanPosts).toBe(1);
  });

  // Stale-permission path: /api/me said the user could scan when the page
  // loaded, but /api/scan now returns 403 — for example because the user was
  // removed from the admin group in the interim. The handler must surface an
  // error toast and restore the button to its idle state so the click isn't
  // silently swallowed and the button isn't stuck in a fake "Scan running"
  // label.
  test('403 on click surfaces an error toast and resets the button', async ({ page }) => {
    await mockMe(page, {
      authEnabled: true,
      email: 'admin@example.com',
      groups: ['provenance-admins'],
      canRunScan: true,
    });

    await page.route('**/api/scan', async (route) => {
      if (route.request().method() === 'POST') {
        await route.fulfill({
          status: 403,
          contentType: 'text/plain',
          body: 'forbidden',
        });
        return;
      }
      await route.continue();
    });

    await page.goto(BASE, { waitUntil: 'networkidle' });
    const btn = page.locator('#btn-scan');
    await btn.click();

    // An error toast appears in #toast-wrap with the "Scan request failed"
    // title and a sub-line that includes the upstream 403.
    const toast = page.locator('#toast-wrap .toast.error', { hasText: 'Scan request failed' });
    await expect(toast).toBeVisible({ timeout: 5_000 });
    await expect(toast).toContainText('403');

    // Button must be back to idle: enabled, label "Run Scan", and never
    // stuck on "Scan running".
    await expect(btn).toBeEnabled();
    await expect(btn).toContainText('Run Scan');
    await expect(btn).not.toContainText(/scan running/i);
  });
});
