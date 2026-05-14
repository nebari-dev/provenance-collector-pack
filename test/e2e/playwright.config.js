// Minimal Playwright config for the dashboard e2e tests. Each spec assumes the
// dashboard is already reachable at DASHBOARD_URL (no webServer / setup is
// owned here — that's the calling workflow's job).

/** @type {import('@playwright/test').PlaywrightTestConfig} */
module.exports = {
  testDir: '.',
  timeout: 30_000,
  expect: { timeout: 10_000 },
  fullyParallel: false,
  workers: 1,
  reporter: [['list'], ['html', { open: 'never', outputFolder: 'playwright-report' }]],
  use: {
    headless: true,
    viewport: { width: 1440, height: 900 },
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },
  projects: [{ name: 'chromium', use: { browserName: 'chromium' } }],
};
