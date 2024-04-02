import { GetCandlesRequest } from '../../../node/grpc/contracts/marketData';
import { useState, useEffect, useCallback } from "react";
import { useAppDispatch, useAppSelector } from '../../../store';
import { setShares } from './spaceSlice';
import { GetTradingSchedulesRequest, GetTradingSchedulesResponse, TradingSchedule } from "../../../node/grpc/contracts/shares";
import { useIpcInvoke, useIpcListen } from "../../hooks";
import { OHLCData, OrderState } from "../../../types";

type GetCandlesResponse = OHLCData[];

export const useGetCandles = (): (req: GetCandlesRequest) => Promise<GetCandlesResponse> => useIpcInvoke("GET_CANDLES");
export const useGetTradingSchedules = (): (req: GetTradingSchedulesRequest) => Promise<GetTradingSchedulesResponse> => useIpcInvoke("GET_TRADING_SCHEDULES");
export const useGetShares = () => useIpcInvoke("GET_SHARES");

// TODO: Нужен хук который сам бы хендлил отписку
const useSubscribeCandles = () => useIpcInvoke("SUBSCRIBE_CANDLES");
const useSubscribeOrders = () => useIpcInvoke("SUBSCRIBE_ORDER");
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

            // TODO: Чтобы избежать лагов графика стоит ограничивать размер candles в N айтемов, в зависимости от размера окна и интервала
            setInitialData(candles);
        } catch (e) {
            setError(e);
        } finally {
            setIsLoading(false);
        }
    };

    const subscribeCandles = async () => {
        await subscribe({
            instrumentId: figiOrInstrumentId,
            interval,
        });
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
    const dispatch = useAppDispatch();
    const shares = useAppSelector(store => store.space.shares)
    const [isLoading, setIsLoading] = useState(false);

    const load = useCallback(async () => {
        try {
            if (isLoading) return;
            setIsLoading(true);
            const response: any = await window.ipc.invoke('GET_SHARES_FROM_STORE');
            dispatch(setShares(response.shares))
        } catch (error) {
            console.error(`get shares error: ${error}`);
            dispatch(setShares([]));
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        load();
    }, [])

    return { shares, isLoading }
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
