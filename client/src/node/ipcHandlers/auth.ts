import { ipcMain } from 'electron';
import { ipcEvents } from '../../ipcEvents';
import storage from '../Storage';
import { authService } from '../grpc/auth';
import { createLogger } from '../logger';
import { retry } from '../utils/retry';

const log = createLogger({ controller: 'auth' });


// Сервак не успевает запускаться, но к этому моменту мы уже пытаемся дернуть этот метод, поэтому нужен ретрай
const hasTokenWithRetry = retry(authService.hasToken, {});

ipcMain.handle(ipcEvents.GET_AUTH_INFO, async (e) => {
    try {
        log.info('GET_AUTH_INFO');
        const isSandbox = await storage.get('isSandbox');

        const { HasToken } = await hasTokenWithRetry({});

        return { isAuthorised: !!HasToken, isSandbox };
    } catch (err) {
        log.error('Не удалось получить данные авторизации', err);
        return Promise.reject('Не удалось получить данные авторизации: ' + err)
    }
});

ipcMain.handle(ipcEvents.PRUNE_TOKENS, async (e) => {
    try {
        return await authService.pruneTokens({});
    } catch (err) {
        return Promise.reject('Не удалось очистить токены ' + err)
    }
});
