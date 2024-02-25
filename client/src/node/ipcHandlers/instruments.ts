import { ipcMain } from "electron";
import { ipcEvents } from "../../ipcEvents";
import { sharesService } from "../grpc/instruments";

ipcMain.handle(ipcEvents.GET_SHARES, async (e, req) => {
    const { instrumentStatus } = req;

    if (!instrumentStatus) return Promise.reject('InstrumentStatus обязательный параметр');


    const res = await new Promise((resolve, reject) => {
        sharesService.getShares({
            instrumentStatus
        }, (e,  Instruments) => {
            if (e) return reject(e);
            resolve(Instruments.instruments)
        });
    });
    return res;
}); 