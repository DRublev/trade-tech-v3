import { ipcMain, safeStorage } from "electron";
import { ipcEvents } from "../../ipcEvents";
import { sharesService } from "../grpc/instruments";
import storage from '../Storage';
import { GetInstrumentsRequest, GetSharesResponse, GetTradingSchedulesRequest, GetTradingSchedulesResponse } from "../grpc/contracts/shares";
import { authService } from "../grpc/auth";

ipcMain.handle(ipcEvents.GET_SHARES, (e, req) => getShares(req))
ipcMain.handle(ipcEvents.GET_TRADING_SCHEDULES, (e, req) => getTradingSchedules(req))

ipcMain.handle(ipcEvents.GET_SHARES_FROM_STORE, async (e) => {
    if (!safeStorage.isEncryptionAvailable()) return Promise.reject("Шифрование не доступно");

    try {
        const shares = await storage.get('shares');
        return Promise.resolve({ shares });
    } catch (error) {
        return Promise.reject(`Не удалось получить данные из сторы: ${error}`)
    }
});

export async function getShares(req: GetInstrumentsRequest): Promise<GetSharesResponse> {
    const { instrumentStatus } = req;

    if (!instrumentStatus) return Promise.reject('InstrumentStatus обязательный параметр');

    const { HasToken } = await authService.hasToken({});
    if (!HasToken) return Promise.reject('Not authorized');

    const res: GetSharesResponse = await new Promise((resolve, reject) => {
        sharesService.getShares({
            instrumentStatus
        }, (error, response) => {
            if (error) return reject(error);
            storage.remove('shares')
            resolve(response)
        });
    });
    storage.save('shares', res.instruments);
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
            if (error) return reject(error);
            resolve(response)
        });
    });

    return res;
}