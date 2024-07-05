import { test as base, _electron, type Page } from "@playwright/test";
import { findLatestBuild, parseElectronApp, type ElectronAppInfo } from 'electron-playwright-helpers';

import type { IIpcRenderer } from '../src/preload';

type WindowPatched = typeof window & { ipc: IIpcRenderer }

class Mocker {
    constructor(private page: Page) {
    }

    async setStub(predicate, stub) {
        // Мокаем window.ipc, так, чтобы при ее вызове с {ipcPredicate}, в аргументах прокидывался также {randomHash}
    }

    private getRandomHash() {

        return ''
    }

}

class ElectronApp {
    private appInfo: ElectronAppInfo;
    constructor() {
        // Ищем последний билд
        // При запуске руками из консоли - в дефолтной директроии "out"
        // При запуске через vs code - относительно корня проекта
        const latestBuild = findLatestBuild(process.env.TESTS ? undefined : "./client/out")
        // parse the directory and find paths and other info
        this.appInfo = parseElectronApp(latestBuild)
    }
    launch() {
        return _electron.launch({
            args: [this.appInfo.main],
            executablePath: this.appInfo.executable
        });
    }
}

type Fixtures = {
    mocker: Mocker;
    electron: ElectronApp;
}

export const test = base.extend<Fixtures>({
    mocker: async ({ page }, use) => {
        const mocker = new Mocker(page);
        await use(mocker);
    },
    electron: async ({ page }, use) => {
        await use(new ElectronApp());
    }
});

export { expect, devices, request, selectors, mergeExpects, mergeTests } from '@playwright/test';