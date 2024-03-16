import { ipcMain, safeStorage } from "electron";
import { ipcEvents } from "../../ipcEvents";
import { sharesService } from "../grpc/instruments";
import storage from '../Storage';
import { GetInstrumentsRequest, GetSharesResponse } from "../../../grpcGW/shares";
import { authService } from "../grpc/auth";

ipcMain.handle(ipcEvents.GET_SHARES, (e, req) => getShares(req))

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
