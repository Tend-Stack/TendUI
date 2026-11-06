const { test, expect } = require('@playwright/test');

test.beforeEach(async ({ context }) => {
  await context.addCookies([{ name: 'tendui_test_htmx', value: '1', url: 'http://127.0.0.1:8089', sameSite: 'Strict' }]);
});

test('consumer-supplied HTMX swaps the edit panel and validation response', async ({ page }) => {
  await page.goto('/editing');
  await page.waitForFunction(() => typeof window.htmx !== 'undefined');
  const editing = page.waitForResponse(r => r.url().includes('/editing?edit=1') && r.request().headers()['hx-request'] === 'true');
  await page.getByRole('link', { name: 'Edit sample profile' }).click();
  expect((await editing).headers()['vary']).toContain('HX-Request');
  await expect(page.getByLabel('Display name')).toBeVisible();
  await expect(page).toHaveURL(/\/editing$/);
  await page.getByLabel('Display name').fill('HTMX visitor');
  await page.getByRole('button', { name: 'Save sample', exact: true }).click();
  await expect(page.getByRole('heading', { name: 'HTMX visitor' })).toBeVisible();
  await page.getByRole('link', { name: 'Edit sample profile' }).click();
  await page.getByLabel('Display name').fill('');
  // This example intentionally returns 200 for invalid HTMX fragments so
  // default response handling swaps the form; native invalid POSTs use 422.
  const invalid = page.waitForResponse(r => r.url().endsWith('/editing') && r.request().method() === 'POST' && r.request().headers()['hx-request'] === 'true');
  await page.getByRole('button', { name: 'Test server validation' }).click();
  expect((await invalid).status()).toBe(200);
  await expect(page.getByLabel('Display name')).toHaveAttribute('aria-invalid', 'true');
  await expect(page.locator('html')).toHaveCount(1);
  await expect(page.locator('#profile-panel')).toHaveCount(1);
});

test('HTMX theme redirect removes preview parameters and applies the cookie', async ({ page }) => {
  await page.goto('/settings?theme=light');
  await page.waitForFunction(() => typeof window.htmx !== 'undefined');
  const toggle = page.waitForResponse(r => r.url().endsWith('/api/theme/toggle'));
  await page.getByRole('button', { name: 'Toggle theme' }).click();
  const response = await toggle;
  expect(response.status()).toBe(200);
  expect(response.headers()['hx-redirect']).toBe('/settings');
  await expect(page).toHaveURL(/\/settings$/);
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark');
});
