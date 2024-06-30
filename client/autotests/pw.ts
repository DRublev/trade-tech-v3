import { test as base, _electron, type Page } from "@playwright/test";
import { findLatestBuild, parseElectronApp, type ElectronAppInfo } from 'electron-playwright-helpers';

import type { IIpcRenderer } from '../src/preload';

const mockServer = {
    setMock: (mock, hash) => {
        // При поступлении запроса с заголовком x-mock в значении {hash} - отдавать {mock}. Но лишь один раз. При след запросе отдавать обычный мок
        // TODO: Делаем запрос на мок-сервер
    }
}
function MockDecorator(mock, ipcPredicate) {
    const randomHash = 'sdasd';
}

type WindowPatched = typeof window & { ipc: IIpcRenderer }

class Mocker {
    constructor(private page: Page) {
        page.addInitScript(() => {
            const originalInvoke = (window as WindowPatched).ipc.invoke;
            window.ipc.invoke = (channel, ...args) => {
                if (channel === 'GET_AUTH_INFO') {
                    return Promise.resolve({ isAuthorised: true, isSandbox: false })
                }
                if (channel === 'GET_ACCOUNT') {
                    return Promise.resolve({ AccountId: "123" })
                }

                return originalInvoke(channel, ...args);
            };
        });
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
        // find the latest build in the out directory
        const latestBuild = findLatestBuild()
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