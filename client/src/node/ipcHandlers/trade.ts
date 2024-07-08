import { BrowserWindow, ipcMain } from 'electron';
import { ipcEvents } from '../../ipcEvents';
import { tradeService } from '../grpc/trade';
import storage from '../Storage';
import type { StrategyEvent } from '../grpc/contracts/trade';

ipcMain.handle(ipcEvents.START_TRADE, async (e, req) => {
    const { instrumentId, strategy } = req;

    if (!instrumentId) return Promise.reject("instrumentId является обязательным параметром");
    if (!strategy) return Promise.reject("strategy является обязательным параметром");

    const response = await tradeService.start({
        InstrumentId: instrumentId,
        Strategy: strategy
    });

    return response
});

ipcMain.handle(ipcEvents.IS_STARTED, async (e, req) => {
    const { instrumentId, strategy } = req;

    if (!instrumentId) return Promise.reject("instrumentId является обязательным параметром");
    if (!strategy) return Promise.reject("strategy является обязательным параметром");

    const response = await tradeService.isStarted({
        InstrumentId: instrumentId,
        Strategy: strategy
    });

    return response
});

ipcMain.handle(ipcEvents.STOP_TRADE, async (e, req) => {
    const { instrumentId, strategy } = req;

    if (!instrumentId) return Promise.reject("instrumentId является обязательным параметром");
    if (!strategy) return Promise.reject("strategy является обязательным параметром");

    const response = await tradeService.stop({
        InstrumentId: instrumentId,
        Strategy: strategy
    });

    return response
});


ipcMain.handle(ipcEvents.CHANGE_STRATEGY_CONFIG, async (e, req) => {
    const { instrumentId, strategy, values } = req;

    if (!instrumentId) return Promise.reject("instrumentId является обязательным параметром");
    if (!strategy) return Promise.reject("strategy является обязательным параметром");
    if (!values) return Promise.reject("values является обязательным параметром");

    try {
        const shares = await storage.get('shares');
        const instrument = shares.find((share: any) => share.uid === instrumentId);
        values.LotSize = instrument.lot;
    } catch (e) {
        console.error(e);
    }


    const response = await tradeService.changeConfig({
        InstrumentId: instrumentId,
        Strategy: strategy,
        Config: values,
    });

    return response;
})

ipcMain.handle(ipcEvents.GET_STRATEGY_CONFIG, async (e, req) => {
    const { instrumentId, strategy } = req;

    if (!instrumentId) return Promise.reject("instrumentId является обязательным параметром");
    if (!strategy) return Promise.reject("strategy является обязательным параметром");

    const response = await tradeService.getConfig({ Strategy: strategy, InstrumentId: instrumentId });
    return response.Config;
});

ipcMain.handle(ipcEvents.SUBSCRIBE_STRATEGY_ACTIVITIES, async (e, req) => {
    const { strategy } = req;
    if (!strategy) return Promise.reject("strategy является обязательным параметром");
    const [win] = BrowserWindow.getAllWindows()
    if (!win) return Promise.reject("не надено ни одного окна приложения");
    const subscription = new Promise((resolve, reject) => {
        const stream = tradeService.subscribeStrategiesEvents({ Strategy: strategy });
        
        stream.on('data', (activity: StrategyEvent) => {
            win.webContents.send(ipcEvents.NEW_STRATEGY_ACTIVITY, activity);
        });
        stream.on('end', () => {
            return resolve(true)
        });
        stream.on('error', (err) => {
            return reject(err);
        });
    });

    return subscription;
});
