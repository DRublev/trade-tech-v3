import { ipcMain } from "electron";
import { ipcEvents } from "../../ipcEvents";
import { sharesService } from "../grpc/instruments";
import storage from '../Storage';
import { GetInstrumentsRequest, GetSharesResponse, GetTradingSchedulesRequest, GetTradingSchedulesResponse } from "../grpc/contracts/shares";
import { authService } from "../grpc/auth";
import { createLogger } from "../logger";

const log = createLogger({ controller: 'instruments' });

ipcMain.handle(ipcEvents.GET_SHARES, (e, req) => getShares(req))
ipcMain.handle(ipcEvents.GET_TRADING_SCHEDULES, (e, req) => getTradingSchedules(req))

ipcMain.handle(ipcEvents.GET_SHARES_FROM_STORE, async (e) => {
    try {
        const shares = await storage.get('shares');
        return Promise.resolve({ shares });
    } catch (error) {
        log.error("Error getting shares from store", error);
        return Promise.reject(`Не удалось получить данные из сторы: ${error}`)
    }
});

export async function getShares(req: GetInstrumentsRequest): Promise<GetSharesResponse> {
    const { instrumentStatus } = req;

    log.info('Getting shares', req);

    if (!instrumentStatus) return Promise.reject('InstrumentStatus обязательный параметр');

    const { HasToken } = await authService.hasToken({});
    if (!HasToken) {
        log.warn("Getting shares without being authorized");
        return Promise.reject('Not authorized');
    }

    const res: GetSharesResponse = await new Promise((resolve, reject) => {
        sharesService.getShares({
            instrumentStatus
        }, (error, response) => {
            if (error) {
                log.error("Error getting shares", error);
                return reject(error);
            }
            storage.remove('shares').then(() => resolve(response));
        });
    });
    await storage.save('shares', res.instruments);
    return res;
}


export async function getTradingSchedules(req: GetTradingSchedulesRequest): Promise<GetTradingSchedulesResponse> {
    const { exchange, from, to } = req;

    if (!from || !to) return Promise.reject('Нет обязательных параметров');

    const res: GetTradingSchedulesResponse = await new Promise((resolve, reject) => {
        sharesService.getTradingSchedules({
            exchange: exchange,
            from: from,
            to: to
        }, (error, response) => {
            if (error) {
                log.error("Error getting schedules", error);
                return reject(error);
            }
            resolve(response)
        });
    });

    return res;
}

ipcMain.handle(ipcEvents.SET_CURRENT_INSTRUMENT, async (e, req) => {
    const { instrumentId } = req;

    if (!instrumentId) return Promise.reject('instrumentId обязательный параметр');

    await storage.save('currentInstrument', instrumentId);

    return Promise.resolve();
});

ipcMain.handle(ipcEvents.GET_CURRENT_INSTRUMENT, async () => {
    const stored = await storage.get('currentInstrument');

    if (!stored) return Promise.reject('Нет сохраненного инструмента');

    return Promise.resolve(stored);
})
