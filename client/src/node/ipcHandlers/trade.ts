import { ipcMain } from 'electron';
import { ipcEvents } from '../../ipcEvents';
import { tradeService } from '../grpc/trade';

ipcMain.handle(ipcEvents.START_TRADE, async (e) => {
    const response = await new Promise(resolve => tradeService.start({
        InstrumentId: 'BBG004730RP0',
        Strategy: 'spread_v0'
    }, (err, res) => {
        if (err) return resolve(false);
        console.log('11 trade', res);

        resolve(true);
    }));
    return response
});
