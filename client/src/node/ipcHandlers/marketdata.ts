import { BrowserWindow, ipcMain, ipcRenderer } from "electron";
import { ipcEvents } from "../../ipcEvents";
import { marketdataService } from "../grpc/marketdata";
import { Quant } from "./types";
import { OHLC, OrderState } from '../grpc/contracts/marketData';
import { OHLCData, OrderState as Order } from "../../types";
import { UTCTimestamp } from "lightweight-charts";
import { ClientReadableStream } from "@grpc/grpc-js";
import { createLogger } from "../logger";

const log = createLogger({ controller: 'marketdata' });

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
    time: candle.time.valueOf() / 1000 as UTCTimestamp,
})
const toOrderState = (orderState: Omit<OrderState, 'time'> & { time: string | Date }): Order => ({
    id: orderState.IdempodentID,
    instrumentId: orderState.InstrumentID,
    price: orderState.PricePerLot,
    status: orderState.ExecutionStatus,
    lotsRequested: orderState.LotsRequested,
    lotsExecuted: orderState.LotsExecuted,
    operationType: orderState.OperationType,
    time: (typeof orderState.time === 'string' ? new Date(orderState.time) : orderState.time).valueOf() / 1000 as UTCTimestamp,
    strategy: orderState.Strategy,
})

ipcMain.handle(ipcEvents.GET_CANDLES, async (e, req) => {
    const { instrumentId, start, end, interval } = req;

    log.info('Get candles');

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
            if (e) {
                log.error("Getting candles error", e);
                return reject(e);
            }
            const { candles } = resp || { candles: [] };
            resolve(candles.map(candleToOhlc))
        });
    });
    return res;
});

// TODO: Вынести бы в класс, но пока пофиг 
type HandleOrderCallback = (order: OrderState, error?: Error) => void;
const subscribers: HandleOrderCallback[] = [];
let stream: ClientReadableStream<OrderState>;
const createStream = () => {
    stream = marketdataService.subscribeOrders({})
    stream.on('data', async (order: OrderState) => {
        await Promise.allSettled(subscribers.map(cb => cb(order)))
    });
    stream.on('end', () => {
        Promise.allSettled(subscribers.map(cb => cb(null, new Error('end of stream'))))

    });
    stream.on('error', (err) => {
        log.error("Error in marketdata stream", err);

        Promise.allSettled(subscribers.map(cb => cb(null, err)))
    });
};

const subscribeForOrderStateChange = (callback: HandleOrderCallback) => {
    subscribers.push(callback)
    if (!stream) {
        createStream();
    }
};

ipcMain.handle(ipcEvents.SUBSCRIBE_ORDER, async (e, req) => {
    const { instrumentId } = req;

    if (!instrumentId) return Promise.reject('InstrumentId обязательный параметр');


    const [win] = BrowserWindow.getAllWindows();
    const res = await new Promise((resolve, reject) => {
        const s = marketdataService.subscribeOrders({})
        s.on('data', async (order: OrderState) => {
            log.info("New order " + JSON.stringify(order))

            if (!order || order.InstrumentID != instrumentId) return;
            try {
                win.webContents.send(ipcEvents.NEW_ORDER, toOrderState(order));
            } catch (e) {
                log.error('Error in marketdata stream', e)
            }
        });
        s.on('end', () => {
            // Promise.allSettled(subscribers.map(cb => cb(null, new Error('end of stream'))))
            return resolve(true)
        });
        s.on('error', (err) => {
            log.error("Error in marketdata stream", err);

            return reject(err);
        });
    });
    // TODO: Сюда бы обработку ошибок (не лог)

    // TODO: Возвращать метод/строку для отписки
    return res;
});

const subscribed: Record<string, boolean> = {};

ipcMain.handle(ipcEvents.UNSUBSCRIBE_CANDLES, async (e, req) => {
    const { instrumentId } = req;
    if (!instrumentId) return Promise.reject('InstrumentId обязательный параметр');
    if (subscribed[instrumentId]) {
        delete subscribed[instrumentId];
    }
    return Promise.resolve(true);
});

ipcMain.handle(ipcEvents.SUBSCRIBE_CANDLES, async (e, req) => {
    const { instrumentId, interval } = req;

    if (!instrumentId) return Promise.reject('InstrumentId обязательный параметр');
    if (!interval) return Promise.reject('interval обязательный параметр');

    const [win] = BrowserWindow.getAllWindows()

    const res = new Promise((resolve, reject) => {
        try {
            // TODO: Хорошо бы это делать в воркере или background процессе
            const stream = marketdataService.subscribeCandles({ instrumentId, interval });
            subscribed[instrumentId] = true;
            stream.on('data', (candle: OHLC) => {
                if (!subscribed[instrumentId]) return;
                win.webContents.send(ipcEvents.NEW_CANDLE, candleToOhlc(candle));
            });
            stream.on('end', () => {
                delete subscribed[instrumentId];
                resolve(true);

            });
            stream.on('error', (err) => {
                log.error("Candles stream error", err);
                reject(err);
            })
            resolve(true);
        } catch (e) {
            log.error("Candles stream creation error", e);
            reject(e);
        }
    });

    // TODO: Возвращать метод/строку для отписки
    return res;
})