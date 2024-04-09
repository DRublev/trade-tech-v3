import { ipcMain, safeStorage } from "electron";
import { ipcEvents } from "../../ipcEvents";
import storage from '../Storage';
import { accountsService } from '../grpc/accounts';
import { authService } from '../grpc/auth';

type Payload = {
    token: string;
    isSandbox?: boolean;
};
ipcMain.handle(ipcEvents.REGISTER, async (e, data: Payload) => {
    if (!safeStorage.isEncryptionAvailable()) return Promise.reject("Шифрование не доступно");

    const { token, isSandbox } = data;
    if (!token) return Promise.reject("token является обязательным параметром");

    await storage.save('isSandbox', isSandbox ? 1 : 0);

    return authService.setToken({ Token: token, IsSandbox: isSandbox });
});

ipcMain.handle(ipcEvents.GET_ACCOUNTS, async (e) => {
    const res = await new Promise((resolve, reject) => {
        accountsService.getAccounts({}, (e, accs) => {
            if (e) return reject(e);
            resolve(accs)
        });
    });
    return res;
});

ipcMain.handle(ipcEvents.SET_ACCOUNT, async (e, data) => {
    if (!data.id) return Promise.reject('id является обязательным параметром');

    storage.save('accountId', data.id);
    return new Promise((resolve, reject) => accountsService.setAccount({ AccountId: data.id }, (e, res) => {
        if (e) return reject(e);
        resolve({});
    }));
});