import { test, expect } from "../pw";


/**
 * @see https://github.com/spaceagetv/electron-playwright-example/blob/main/e2e-tests/main.spec.ts
 */


test("Проверяем, что отображается график", async ({ electron }) => {
  const app = await electron.launch();

  const win = await app.firstWindow();
  await win.waitForSelector('#chart-container', { timeout: 1000 })
  const chart = await win.$('#chart-container');

  await expect(chart).toBeTruthy();
});
