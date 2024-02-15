import { ipcMain, safeStorage } from 'electron';
import { ipcEvents } from '../../ipcEvents';
import { testService } from '../grpc/test';
import storage from '../Storage';

ipcMain.handle(ipcEvents.TEST_HELLO, async (e, data) => {
    console.log(ipcEvents.TEST_HELLO, e, data);
    testService.ping({ content: 'from client' }, (e, res) => {
        console.log('e ', e);
        console.log('res ', res)
    });
});

type Payload = {
    token: string;
    isSandbox?: boolean;
};
ipcMain.handle(ipcEvents.REGISTER, async (e, data: Payload) => {
    if (!safeStorage.isEncryptionAvailable()) return Promise.reject("Шифрование не доступно");

    const { token, isSandbox } = data;
    if (!token) return Promise.reject("Token is required!");

    const encryptedToken = safeStorage.encryptString(token);
    await storage.save('isSandbox', isSandbox ? 1 : 0);
    await storage.save('token', encryptedToken);

});