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

ipcMain.handle(ipcEvents.IS_STARTED, async (e, req) => {
    const { instrumentId } = req;

    if (!instrumentId) return Promise.reject("instrumentId является обязательным параметром");

    const response = await tradeService.isStarted({
        InstrumentId: instrumentId,
        // TODO: Сделать эндпоинт, который бы отдавал с бека список стратегий с некоторой инфой, а на фронте его показывать как опции
        Strategy: 'spread_v0'
    });

    return response
});

ipcMain.handle(ipcEvents.STOP_TRADE, async (e, req) => {
    const { instrumentId } = req;

    if (!instrumentId) return Promise.reject("instrumentId является обязательным параметром");

    const response = await tradeService.stop({
        InstrumentId: instrumentId,
        // TODO: Сделать эндпоинт, который бы отдавал с бека список стратегий с некоторой инфой, а на фронте его показывать как опции
        Strategy: 'spread_v0'
    });

    return response
});


ipcMain.handle(ipcEvents.CHANGE_STRATEGY_CONFIG, async (e, req) => {
    const { instrumentId, strategy, values } = req;

    if (!instrumentId) return Promise.reject("instrumentId является обязательным параметром");
    if (!strategy) return Promise.reject("strategy является обязательным параметром");
    if (!values) return Promise.reject("values является обязательным параметром");

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
