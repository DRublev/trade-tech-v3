import { ipcMain } from 'electron';
import { ipcEvents } from '../../ipcEvents';
import storage from '../Storage';
import { authService } from '../grpc/auth';

ipcMain.handle(ipcEvents.GET_AUTH_INFO, async (e) => {
    try {
        const isSandbox = await storage.get('isSandbox');
        const account = await storage.get('accountId');

        const { HasToken } = await authService.hasToken({});

        return { isAuthorised: !!HasToken, isSandbox, account };
    } catch (err) {
        return Promise.reject('Не удалось получить данные из сторы: ' + err)
    }
});

ipcMain.handle(ipcEvents.PRUNE_TOKENS, async (e) => {
    try {
        return await authService.pruneTokens({});
    } catch (err) {
        return Promise.reject('Не удалось очистить токены ' + err)
    }
});
