const { test, expect } = require('@playwright/test');
const AxeBuilder = require('@axe-core/playwright').default;
const fs = require('node:fs');
const path = require('node:path');
const source = fs.readFileSync(path.resolve(__dirname, '../../examples/site/catalog.go'), 'utf8');
const names = [...source.matchAll(/^\s*\{"([A-Z][A-Za-z]+)",/gm)].map(m => m[1].toLowerCase());

test('catalog test inventory covers all 36 components', async () => {
  expect(new Set(names).size).toBe(36);
});

for (const mode of ['light', 'dark']) {
  test(`all component details: ${mode} accessibility rules and mobile overflow`, async ({ page }) => {
    test.setTimeout(180000);
    await page.setViewportSize({ width: 390, height: 844 });
    for (const name of names) {
      await page.goto(`/catalog/${name}?theme=${mode}`);
      await expect(page.locator('h1')).toHaveCount(1);
      const scan = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag21aa']).analyze();
      expect(scan.violations.map(v => ({ id: v.id, nodes: v.nodes.map(n => n.target) })), name).toEqual([]);
      const widths = await page.evaluate(() => ({ page: document.documentElement.scrollWidth, viewport: innerWidth }));
      expect(widths.page, name).toBeLessThanOrEqual(widths.viewport + 1);
    }
  });
}

for (const width of [320, 768, 1920]) {
  test(`application layouts remain contained at ${width}px`, async ({ page }) => {
    await page.setViewportSize({ width, height: 1080 });
    for (const route of ['/playground', '/dashboard', '/workspace', '/forms', '/patterns', '/settings']) {
      await page.goto(route);
      expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1), route).toBe(true);
    }
  });
}

test('keyboard focus is visible and skip link reaches main content', async ({ page }) => {
  await page.goto('/forms');
  await page.keyboard.press('Tab');
  await expect(page.getByRole('link', { name: 'Skip to content' })).toBeFocused();
  const ring = await page.locator(':focus').evaluate(e => getComputedStyle(e).outlineStyle);
  expect(ring).not.toBe('none');
  await page.keyboard.press('Enter');
  await expect(page).toHaveURL(/#main$/);
});

test('reduced motion and explicit off stop animation without hiding content', async ({ page }) => {
  await page.emulateMedia({ reducedMotion: 'reduce' });
  await page.goto('/motion');
  expect(await page.locator('[data-enter], .skeleton-shapes').evaluateAll(nodes => nodes.every(n => getComputedStyle(n).animationName === 'none'))).toBe(true);
  await page.emulateMedia({ reducedMotion: 'no-preference' });
  await page.goto('/motion?motion=off');
  expect(await page.locator('[data-enter], .skeleton-shapes').evaluateAll(nodes => nodes.every(n => getComputedStyle(n).animationName === 'none'))).toBe(true);
  await expect(page.locator('main h1')).toBeVisible();
});

test('forced colors retain native controls and readable content', async ({ page, browserName }) => {
  test.skip(browserName === 'webkit', 'WebKit does not support forced-colors emulation; review on supporting platforms.');
  await page.emulateMedia({ forcedColors: 'active' });
  await page.goto('/playground');
  await expect(page.getByLabel('A native checkbox', { exact: true })).toBeVisible();
  await page.getByLabel('A native checkbox', { exact: true }).check();
  await expect(page.getByLabel('A native checkbox', { exact: true })).toBeChecked();
});

test('print palette overrides a dark page and produces a review PDF', async ({ page, browserName }, testInfo) => {
  await page.goto('/dashboard?theme=dark');
  await page.emulateMedia({ media: 'print' });
  const colors = await page.locator('html').evaluate(e => {
    const s = getComputedStyle(e); return [s.getPropertyValue('--gothui-background').trim(), s.getPropertyValue('--gothui-text').trim()];
  });
  expect(colors).toEqual(['#ffffff', '#000000']);
  if (browserName === 'chromium') {
    const pdf = await page.pdf({ format: 'A4', printBackground: false });
    await testInfo.attach('dashboard-print', { body: pdf, contentType: 'application/pdf' });
  }
});
