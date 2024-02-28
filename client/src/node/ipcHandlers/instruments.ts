import { ipcMain, safeStorage } from "electron";
import { ipcEvents } from "../../ipcEvents";
import { sharesService } from "../grpc/instruments";
import storage from '../Storage';
import { GetInstrumentsRequest, GetSharesResponse } from "../../../grpcGW/shares";

ipcMain.handle(ipcEvents.GET_SHARES, (e, req) => getShares(req))

ipcMain.handle(ipcEvents.GET_SHARES_FROM_STORE, async (e) => {
    if (!safeStorage.isEncryptionAvailable()) return Promise.reject("Шифрование не доступно");

    try {
        const shares = await storage.get('shares');
        return Promise.resolve({ shares });
    } catch (err) {
        return Promise.reject('Не удалось получить данные из сторы: ' + err)
    }
});

export async function getShares(req: GetInstrumentsRequest): Promise<GetSharesResponse> {
    const { instrumentStatus } = req;

    if (!instrumentStatus) return Promise.reject('InstrumentStatus обязательный параметр');

    const res: any = await new Promise((resolve, reject) => {
        sharesService.getShares({
            instrumentStatus
        }, (e, resp) => {
            if (e) return reject(e);
            storage.remove('shares')
            resolve(resp.instruments)
        });
    });
    storage.save('shares', res);
    return res;
}
