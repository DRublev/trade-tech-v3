import { BrowserWindow, ipcMain, ipcRenderer } from "electron";
import { ipcEvents } from "../../ipcEvents";
import { marketdataService } from "../grpc/marketdata";
import { Quant } from "./types";
import { OHLC } from "../../../grpcGW/marketData";
import { OHLCData } from "../../types";

const nanoPrecision = 1_000_000_000;
const quantToNumber = (q: Quant): number => {
    return Number(q.units + (q.nano / nanoPrecision));
}

const candleToOhlc = (candle: OHLC): OHLCData => ({
    open: quantToNumber(candle.open),
    high: quantToNumber(candle.high),
    low: quantToNumber(candle.low),
    close: quantToNumber(candle.close),
    volume: candle.volume,
    date: candle.time,
})

ipcMain.handle(ipcEvents.GET_CANDLES, async (e, req) => {
    const { instrumentId, start, end, interval } = req;

    if (!instrumentId) return Promise.reject('InstrumentId обязательный параметр');
    if (!start) return Promise.reject('start обязательный параметр');
    if (!interval) return Promise.reject('interval обязательный параметр');

    const res = await new Promise((resolve, reject) => {
        marketdataService.getCandles({
            instrumentId,
            interval,
            start,
            end: end || Date.now(),
        }, (e, resp) => {
            if (e) return reject(e);
            const { candles } = resp || { candles: [] };
            resolve(candles.map(candleToOhlc))
        });
    });
    return res;
});

ipcMain.handle(ipcEvents.SUBSCRIBE_CANDLES, async (e, req) => {
    const { instrumentId, interval } = req;

    if (!instrumentId) return Promise.reject('InstrumentId обязательный параметр');
    if (!interval) return Promise.reject('interval обязательный параметр');

    const [win] = BrowserWindow.getAllWindows()

    const res = new Promise(async (resolve, reject) => {
        try {
            // TODO: Хорошо бы это делать в воркере или background процессе
            const stream = marketdataService.subscribeCandles({ instrumentId, interval })
            stream.on('data', candle => {
                console.log("53 marketdata", candle);
                win.webContents.send(ipcEvents.NEW_CANDLE, {
                    ...candle,
                    open: quantToNumber(candle.open),
                    high: quantToNumber(candle.high),
                    low: quantToNumber(candle.low),
                    close: quantToNumber(candle.close),
                    // TODO: Пофиксить на беке
                    date: candle.time
                });
            });
            stream.on('end', () => {
                resolve(true);
            });
            stream.on('error', (err) => {
                console.log("60 marketdata", err);
                reject(err);
            })
            resolve(true);
        } catch (e) {
            reject(e);
        }
    });

    // TODO: Возвращать метод/строку для отписки
    return res;
})