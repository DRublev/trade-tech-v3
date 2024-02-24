import { ipcMain } from "electron";
import { ipcEvents } from "../../ipcEvents";
import { instrumentsService } from "../grpc/instruments";

ipcMain.handle(ipcEvents.GET_SHARES, async (e, req) => {
    const { instrumentStatus } = req;

    if (!instrumentStatus) return Promise.reject('InstrumentId обязательный параметр');


    const res = await new Promise((resolve, reject) => {
        instrumentsService.getShares({
            instrumentStatus
        }, (e, { Instruments }) => {
            if (e) return reject(e);
            resolve(Instruments)
        });
    });
    return res;
}); 