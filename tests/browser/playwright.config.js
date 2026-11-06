const { defineConfig } = require('@playwright/test');
module.exports = defineConfig({
  testDir: '.',
  testMatch: '*.spec.js',
  outputDir: '../results/browser/artifacts',
  reporter: [['line'], ['json', { outputFile: '../results/browser/results.json' }]],
  fullyParallel: true,
  workers: 3,
  retries: 0,
  timeout: 45000,
  expect: { timeout: 8000 },
  use: { baseURL: 'http://127.0.0.1:8089', trace: 'retain-on-failure', screenshot: 'only-on-failure' },
  projects: [
    { name: 'chromium', use: { browserName: 'chromium' } },
    { name: 'firefox', use: { browserName: 'firefox' } },
    { name: 'webkit', use: { browserName: 'webkit' } },
  ],
  webServer: {
    command: '../results/browser/fixture -htmx node_modules/htmx.org/dist/htmx.min.js',
    url: 'http://127.0.0.1:8089/catalog',
    reuseExistingServer: false,
    timeout: 15000,
  },
});
