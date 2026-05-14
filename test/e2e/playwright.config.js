// Minimal Playwright config for the dashboard e2e tests. Each spec assumes the
// dashboard is already reachable at DASHBOARD_URL (no webServer / setup is
// owned here — that's the calling workflow's job).
//
// Video and trace are recorded on every run (not just on failure) so the CI
// artifacts always include a visual record of what the timeline-switch test
// actually did. Reviewers can scrub the .webm file or open the trace in
// Playwright's trace viewer (`npx playwright show-trace trace.zip`).

/** @type {import('@playwright/test').PlaywrightTestConfig} */
module.exports = {
  testDir: '.',
  timeout: 30_000,
  expect: { timeout: 10_000 },
  fullyParallel: false,
  workers: 1,
  outputDir: 'test-results',
  reporter: [['list'], ['html', { open: 'never', outputFolder: 'playwright-report' }]],
  use: {
    headless: true,
    viewport: { width: 1440, height: 900 },
    trace: 'on',
    video: 'on',
    screenshot: 'only-on-failure',
  },
  projects: [{ name: 'chromium', use: { browserName: 'chromium' } }],
};
