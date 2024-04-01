import { GetCandlesRequest } from '../../../node/grpc/contracts/marketData';
import { useIpcInoke, useIpcListen } from "../../hooks";
import { OHLCData, OrderState } from "../../../types";
import { useState, useEffect, useCallback } from "react";
import { GetTradingSchedulesRequest, GetTradingSchedulesResponse, TradingSchedule } from "../../../node/grpc/contracts/shares";

type GetCandlesResponse = OHLCData[];

export const useGetCandles = (): (req: GetCandlesRequest) => Promise<GetCandlesResponse> => useIpcInoke("GET_CANDLES");
export const useGetTradingSchedules = (): (req: GetTradingSchedulesRequest) => Promise<GetTradingSchedulesResponse> => useIpcInoke("GET_TRADING_SCHEDULES");
export const useGetShares = () => useIpcInoke("GET_SHARES");

// TODO: Нужен хук который сам бы хендлил отписку
const useSubscribeCandles = () => useIpcInoke("SUBSCRIBE_CANDLES");
const useSubscribeOrders = () => useIpcInoke("SUBSCRIBE_ORDER");
const useListenCandles = () => useIpcListen("NEW_CANDLE");
const useListenOrders = () => useIpcListen("NEW_ORDER");


type OnOrderCallback = (d: OrderState) => void;

export const useOrders = (callback: OnOrderCallback | OnOrderCallback[], figiOrInstrumentId: string) => {
    const subscribe = useSubscribeOrders();
    const [registerOrderCb, unregisterOrderCb] = useListenOrders();

    const subscribeOrders = async () => {
        await subscribe({
            instrumentId: figiOrInstrumentId,
        });
    }

    const handleNewOrder = useCallback((e: Event, order: OrderState) => {
        if (!order) return;

        if (Array.isArray(callback)) {
            Promise.all(callback.map(cb => () => cb(order)))
        } else {
            callback(order);
        }
    }, []);

    const unsubscribe = () => {
        unregisterOrderCb(handleNewOrder);
    }

    useEffect(() => {
        registerOrderCb(handleNewOrder);
    }, [handleNewOrder]);

    useEffect(() => {
        subscribeOrders();

        return unsubscribe;
    }, [figiOrInstrumentId]);

    return unsubscribe;
}

export const useCandles = (onNewCandle: (d: OHLCData) => void, figiOrInstrumentId: string, interval = 1) => {
    const getCandles = useGetCandles();
    const subscribe = useSubscribeCandles();

    const [registerCandleCb, unregisterCandleCb] = useListenCandles();

    const [initialData, setInitialData] = useState<OHLCData[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);

    const now = new Date();
    // TODO: Рассчитывать на основании интервала
    const dayAgoDate = new Date(new Date(now).setDate(now.getDate() - 1))

    const getInitialCandels = async () => {
        try {
            setIsLoading(true);

            const candles = await getCandles({
                instrumentId: figiOrInstrumentId,
                interval,
                start: dayAgoDate,
                end: now,
            });

            console.log("39 hooks", new Date(candles[candles.length - 1].time));

            // TODO: Чтобы избежать лагов графика стоит ограничивать размер candles в N айтемов, в зависимости от размера окна и интервала
            setInitialData(candles);
        } catch (e) {
            setError(e);
        } finally {
            setIsLoading(false);
        }
    };

    const subscribeCandles = async () => {
        const res = await subscribe({
            instrumentId: figiOrInstrumentId,
            interval,
        });
        console.log("49 hooks", res);
    }

    const handleNewCandle = useCallback((e: Event, candle: OHLCData) => {
        if (!candle || !candle.time) return;

        onNewCandle(candle);
    }, []);

    useEffect(() => {
        registerCandleCb(handleNewCandle);
    }, [handleNewCandle]);

    useEffect(() => {
        getInitialCandels();
        subscribeCandles();

        return () => {
            unregisterCandleCb(handleNewCandle);
        }
    }, [figiOrInstrumentId]);

    return { initialData, isLoading, error };
}

export const useSharesFromStore = () => {
    const [sharesFromStore, setShares] = useState([]);
    const [isLoading, setIsLoading] = useState(false);

    const load = useCallback(async () => {
        try {
            if (isLoading) return;
            setIsLoading(true);
            const response: any = await window.ipc.invoke('GET_SHARES_FROM_STORE');
            setShares(response.shares)
        } catch (error) {
            console.error(`get shares error: ${error}`);
            setShares([]);
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        load();
    }, [])

    return { sharesFromStore, isLoading }
};

export const getTodaysSchedules = (): TradingSchedule[] => {
    const getSchedules = useGetTradingSchedules();
    const [schedules, setSchedules] = useState([]);
    const [isLoading, setIsLoading] = useState(false);

    const load = useCallback(async () => {
        const now = new Date();

        try {
            if (isLoading) return;
            setIsLoading(true);
            const response: any = (await getSchedules({ exchange: "", from: now, to: now })).exchanges;
            setSchedules(response)
        } catch (error) {
            console.error(`get shares error: ${error}`);
            setSchedules([]);
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        load();
    }, [])

    return schedules
};


// "BBG004730RP0" /* GAZP */
// "4c466956-d2ce-4a95-abb4-17947a65f18a" TGLD
// "BBG004730ZJ9" /* VTBR */
// "BBG004PYF2N3" /* POLY */
let instrumentId = "BBG004PYF2N3" /* POLY */;

export const useCurrentInstrumentId = (): [string, (c: string) => void] => {
    const set = (candidate: string) => {
        if (!candidate) throw new Error('candidate is required');

        instrumentId = candidate;
    };

    return [instrumentId, set];
};
