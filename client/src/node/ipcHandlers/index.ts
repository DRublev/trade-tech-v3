import { ipcMain, safeStorage } from 'electron';
import { ipcEvents } from '../../ipcEvents';
import { testService } from '../grpc/test';
import storage from '../Storage';
import { accountsService } from '../grpc/accounts';

ipcMain.handle(ipcEvents.TEST_HELLO, async (e, data) => {
    console.log(ipcEvents.TEST_HELLO, e, data);
    testService.ping({ content: 'from client' }, (e, res) => {
        console.log('e ', e);
        console.log('res ', res)
    });
});

ipcMain.handle(ipcEvents.GET_AUTH_INFO, async (e) => {
    if (!safeStorage.isEncryptionAvailable()) return Promise.reject("Шифрование не доступно");

    try {
        const token = await storage.get('token');
        const isSandbox = await storage.get('isSandbox');
        const accountId = await storage.get('accountId');

        return { isAuthorised: !!token, isSandbox, accountId };
    } catch (err) {
        return Promise.reject('Не удалось получить данные из сторы: ' + err)
    }
});

type Payload = {
    token: string;
    isSandbox?: boolean;
};
ipcMain.handle(ipcEvents.REGISTER, async (e, data: Payload) => {
    if (!safeStorage.isEncryptionAvailable()) return Promise.reject("Шифрование не доступно");

    const { token, isSandbox } = data;
    if (!token) return Promise.reject("Токен является обязательным параметром");

    // TODO: Шифровать и хранить в go
    const encryptedToken = safeStorage.encryptString(token);
    await storage.save('isSandbox', isSandbox ? 1 : 0);
    await storage.save('token', encryptedToken);
});

ipcMain.handle(ipcEvents.GET_ACCOUNTS, async (e) => {
    console.log("31 index");

    const res = await new Promise(resolve => {
        accountsService.getAccounts({}, (e, accs) => {
            console.log("32 index", e, accs);
            resolve(accs)
        });
    });
    console.log("39 index", res);

    return [];
});
