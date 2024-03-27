import { ipcMain } from 'electron';
import { ipcEvents } from '../../ipcEvents';
import { tradeService } from '../grpc/trade';

ipcMain.handle(ipcEvents.START_TRADE, async (e, req) => {
    const { instrumentId } = req;

    if (!instrumentId) return Promise.reject("instrumentId является обязательным параметром");

    const response = await tradeService.start({
        InstrumentId: instrumentId,
        // TODO: Сделать эндпоинт, который бы отдавал с бека список стратегий с некоторой инфой, а на фронте его показывать как опции
        Strategy: 'spread_v0'
    });

    return response
});
