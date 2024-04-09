import { ipcMain, safeStorage } from 'electron';
import { ipcEvents } from '../../ipcEvents';
import storage from '../Storage';
import { authService } from '../grpc/auth';
import { createLogger } from '../logger';

const log = createLogger({ controller: 'auth' });

ipcMain.handle(ipcEvents.GET_AUTH_INFO, async (e) => {
    if (!safeStorage.isEncryptionAvailable()) return Promise.reject("Шифрование не доступно");

    try {
        log.trace('GET_AUTH_INFO');
        const isSandbox = await storage.get('isSandbox');
        const account = await storage.get('accountId');

        const { HasToken } = await authService.hasToken({});

        return { isAuthorised: !!HasToken, isSandbox, account };
    } catch (err) {
        log.error('Не удалось получить данные авторизации', err);
        return Promise.reject('Не удалось получить данные авторизации: ' + err)
    }
});
