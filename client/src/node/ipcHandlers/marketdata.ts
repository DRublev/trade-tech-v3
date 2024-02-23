import { ipcMain } from "electron";
import { ipcEvents } from "../../ipcEvents";
import { marketdataService } from "../grpc/marketdata";
import { Quant } from "./types";
import { GetCandlesResponse_OHLC } from "../../../grpcGW/marketData";
import { OHLCData } from "../../types";

const nanoPrecision = 1_000_000_000;
const quantToNumber = (q: Quant): number => {
    return Number(q.units + (q.nano / nanoPrecision));
}

const candleToOhlc = (candle: GetCandlesResponse_OHLC): OHLCData => ({
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
        }, (e, { candles }) => {
            if (e) return reject(e);
            resolve(candles.map(candleToOhlc))
        });
    });
    return res;
});