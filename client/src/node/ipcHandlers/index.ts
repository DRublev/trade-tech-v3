import { ipcMain, safeStorage } from 'electron';
import { ipcEvents } from '../../ipcEvents';
import storage from '../Storage';

import './register';

ipcMain.handle(ipcEvents.TEST_HELLO, async (e, data) => {
    console.log(ipcEvents.TEST_HELLO, e, data);
});

ipcMain.handle(ipcEvents.GET_AUTH_INFO, async (e) => {
    if (!safeStorage.isEncryptionAvailable()) return Promise.reject("Шифрование не доступно");

    try {
        const isAuthorised = await storage.get('isAuthorised');
        const isSandbox = await storage.get('isSandbox');
        const account = await storage.get('accountId');

        return { isAuthorised: !!isAuthorised, isSandbox, account };
    } catch (err) {
        return Promise.reject('Не удалось получить данные из сторы: ' + err)
    }
});

