const { test, expect } = require('@playwright/test');

test.use({ javaScriptEnabled: false });

test('native profile save, cancel, and validation', async ({ page }) => {
  await page.goto('/editing');
  await page.getByRole('link', { name: 'Edit sample profile' }).click();
  await page.getByLabel('Display name').fill('Native visitor');
  await page.getByLabel('Introduction').fill('No scripts needed.');
  await page.getByRole('button', { name: 'Save sample', exact: true }).click();
  await expect(page.getByRole('heading', { name: 'Native visitor' })).toBeVisible();
  await page.reload();
  await expect(page.getByRole('heading', { name: 'Native visitor' })).toBeVisible();
  await page.getByRole('link', { name: 'Edit sample profile' }).click();
  await page.getByLabel('Display name').fill('Discard this');
  await page.getByRole('link', { name: 'Cancel', exact: true }).click();
  await expect(page.getByRole('heading', { name: 'Native visitor' })).toBeVisible();
  await page.getByRole('link', { name: 'Edit sample profile' }).click();
  await page.getByLabel('Display name').fill('');
  await page.getByRole('button', { name: 'Test server validation' }).click();
  await expect(page.getByLabel('Display name')).toHaveAttribute('aria-invalid', 'true');
  await page.locator('a[href="#profile-name"]').click();
  await expect(page.getByLabel('Display name')).toBeFocused();
});

test('appearance preferences persist across pages', async ({ page }) => {
  await page.goto('/settings');
  await page.getByRole('combobox', { name: 'Color mode', exact: true }).selectOption('dark');
  await page.getByRole('button', { name: 'Save preferences' }).click();
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark');
  await page.goto('/catalog');
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark');
  await page.getByRole('button', { name: 'Toggle theme' }).click();
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'light');
  await expect(page).toHaveURL(/\/catalog$/);
});

test('keyboard opens and closes native popover with focus return', async ({ page }) => {
  await page.goto('/interactions');
  const opener = page.getByRole('button', { name: 'Show project details' });
  await opener.focus();
  await page.keyboard.press('Enter');
  await expect(page.locator('#project-help')).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(page.locator('#project-help')).toBeHidden();
  await expect(opener).toBeFocused();
});

test('native modal keeps the background inert and returns focus on Escape', async ({ page }) => {
  await page.goto('/interactions');
  const opener = page.getByRole('button', { name: 'Review demo request', exact: true });
  await opener.focus();
  await page.keyboard.press('Enter');
  const dialog = page.locator('#approval-dialog');
  await expect(dialog).toBeVisible();
  await expect(page.locator('#approval-dialog:modal')).toHaveCount(1);
  const close = dialog.getByRole('button', { name: 'Close', exact: true });
  await expect(close).toBeFocused();
  await page.keyboard.press('Tab');
  await expect(dialog.getByRole('button', { name: 'Approve demo request' })).toBeFocused();
  await page.keyboard.press('Shift+Tab');
  await expect(close).toBeFocused();
  // Native modality prevents focus moving into the background. Browser chrome
  // remains reachable; do not impose a JavaScript-style circular focus trap.
  await opener.focus();
  await expect(close).toBeFocused();
  await page.keyboard.press('Escape');
  await expect(dialog).toBeHidden();
  await expect(opener).toBeFocused();
});

test('ordinary-page confirmation remains functional without scripting', async ({ page }) => {
  await page.goto('/interactions');
  await page.getByRole('link', { name: 'Use the ordinary-page confirmation' }).click();
  await page.locator('#inline-confirm').getByRole('button', { name: 'Approve demo request' }).click();
  await expect(page).toHaveURL(/confirmed=1/);
});

test('catalog search and downloadable source work natively', async ({ page }) => {
  await page.goto('/catalog');
  await page.getByRole('searchbox', { name: 'Search components' }).fill('DataTable');
  await page.getByRole('button', { name: 'Find components' }).click();
  await page.getByRole('heading', { name: 'DataTable', exact: true }).getByRole('link').click();
  const downloadEvent = page.waitForEvent('download');
  await page.getByRole('link', { name: 'Download example.templ' }).click();
  expect((await downloadEvent).suggestedFilename()).toBe('datatable_example.templ');
});


test('grouped menu and playground scenarios work without scripts', async ({ page }) => {
  await page.goto('/catalog');
  const menu = page.locator('.demo-menu > summary');
  await menu.focus();
  await page.keyboard.press('Enter');
  await expect(page.getByRole('link', { name: 'Form lab', exact: true })).toBeVisible();
  await page.getByRole('link', { name: 'Form lab', exact: true }).click();
  await expect(page).toHaveURL(/\/forms$/);
  await page.getByRole('navigation', { name: 'Demo pages', exact: true }).getByRole('link', { name: 'Playground', exact: true }).click();
  await page.getByRole('navigation', { name: 'Try a scenario' }).getByRole('link', { name: 'Validation errors' }).click();
  await expect(page.getByLabel('Example text field', { exact: true })).toHaveAttribute('aria-invalid', 'true');
  await page.getByLabel('Button label', { exact: true }).fill('Send sample');
  await page.getByRole('combobox', { name: /^Control state/ }).selectOption('readonly');
  await page.getByRole('button', { name: 'Apply and preview' }).click();
  await expect(page).toHaveURL(/#component-preview$/);
  await expect(page.getByLabel('Example text field', { exact: true })).toHaveAttribute('readonly', '');
  await page.getByRole('button', { name: 'Send sample', exact: true }).click();
  await expect(page.getByRole('status')).toContainText('Button submitted successfully');
});
