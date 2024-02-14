import { ipcMain } from 'electron';
import { ipcEvents } from '../ipcEvents';
import { testService } from '../grpc/test';

ipcMain.handle(ipcEvents.TEST_HELLO, async (e, data) => {
    console.log(ipcEvents.TEST_HELLO, e, data);
    testService.ping({ content: 'from client' }, (e, res) => {
        console.log('e ', e);
        console.log('res ', res)
    });
});