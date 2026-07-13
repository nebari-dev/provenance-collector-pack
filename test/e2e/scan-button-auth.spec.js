// Run Scan button visibility based on the authenticated user's groups.
//
// The button is only rendered when `/api/me` returns `canRunScan: true`. This
// spec exercises that frontend gate by mocking the `/api/me` response — no real
// OIDC issuer involved. Server-side enforcement (the bearer → userinfo →
// admin-group check) is covered by internal/dashboard/auth_test.go and
// internal/dashboard/internal_server_test.go.
//
// DASHBOARD_URL points at the React SPA (served by Vite in CI), which proxies
// /api/* to the dashboard. Auth is injected via window.__PW_E2E_AUTH__ so the
// SPA skips the real Keycloak login (honored only outside production builds).

const { test, expect } = require('@playwright/test');

const BASE = process.env.DASHBOARD_URL || 'http://localhost:5173';

async function injectAuth(page) {
  await page.addInitScript(() => {
    window.__PW_E2E_AUTH__ = {
      authenticated: true,
      token: 'e2e-token',
      idTokenParsed: { name: 'E2E Admin', email: 'admin@example.com', preferred_username: 'admin' },
    };
  });
}

// mockMe sets up a route handler that returns the given /api/me payload. Must
// be called BEFORE page.goto() so the response is in place when the SPA fetches
// it on load.
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
  test.beforeEach(async ({ page }) => {
    await injectAuth(page);
  });

  test('absent when auth is disabled', async ({ page }) => {
    await mockMe(page, { authEnabled: false, canRunScan: false, features: {} });
    await page.goto(BASE, { waitUntil: 'networkidle' });

    await expect(page.getByTestId('run-scan')).toHaveCount(0);
  });

  test('absent when user is authenticated but lacks an admin group', async ({ page }) => {
    await mockMe(page, {
      authEnabled: true,
      email: 'alice@example.com',
      groups: ['platform-users'],
      canRunScan: false,
      features: {},
    });
    await page.goto(BASE, { waitUntil: 'networkidle' });

    await expect(page.getByTestId('run-scan')).toHaveCount(0);
  });

  test('visible when user is in an admin group', async ({ page }) => {
    await mockMe(page, {
      authEnabled: true,
      email: 'admin@example.com',
      groups: ['provenance-admins'],
      canRunScan: true,
      features: {},
    });
    await page.goto(BASE, { waitUntil: 'networkidle' });

    const btn = page.getByTestId('run-scan');
    await expect(btn).toBeVisible();
    await expect(btn).toBeEnabled();
  });

  test('clicking the button posts to /api/scan', async ({ page }) => {
    await mockMe(page, {
      authEnabled: true,
      email: 'admin@example.com',
      groups: ['provenance-admins'],
      canRunScan: true,
      features: {},
    });

    let scanPosts = 0;
    await page.route('**/api/scan', async (route) => {
      if (route.request().method() === 'POST') {
        scanPosts++;
        await route.fulfill({
          status: 202,
          contentType: 'application/json',
          body: JSON.stringify({ jobName: 'manual-scan-spec-stub', namespace: 'default' }),
        });
        return;
      }
      await route.continue();
    });

    await page.goto(BASE, { waitUntil: 'networkidle' });
    await page.getByTestId('run-scan').click();

    // While the post is in flight and the poller runs, the button label flips
    // to "Scan running". Seeing that confirms the POST went out.
    await expect(page.getByTestId('run-scan')).toContainText(/scan running/i, { timeout: 5_000 });
    expect(scanPosts).toBe(1);
  });

  // Stale-permission path: /api/me said the user could scan when the page
  // loaded, but /api/scan now returns 403. The hook must surface an error toast
  // and restore the button to its idle state.
  test('403 on click surfaces an error toast and resets the button', async ({ page }) => {
    await mockMe(page, {
      authEnabled: true,
      email: 'admin@example.com',
      groups: ['provenance-admins'],
      canRunScan: true,
      features: {},
    });

    await page.route('**/api/scan', async (route) => {
      if (route.request().method() === 'POST') {
        await route.fulfill({ status: 403, contentType: 'text/plain', body: 'forbidden' });
        return;
      }
      await route.continue();
    });

    await page.goto(BASE, { waitUntil: 'networkidle' });
    const btn = page.getByTestId('run-scan');
    await btn.click();

    const toast = page.getByTestId('toast-error').filter({ hasText: 'Scan request failed' });
    await expect(toast).toBeVisible({ timeout: 5_000 });
    await expect(toast).toContainText('403');

    await expect(btn).toBeEnabled();
    await expect(btn).toContainText('Run Scan');
    await expect(btn).not.toContainText(/scan running/i);
  });
});
