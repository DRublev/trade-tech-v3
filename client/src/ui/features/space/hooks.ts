import { GetCandlesRequest } from '.././contracts/marketData';
import { useIpcInvoke, useIpcListen } from "../../hooks";
import { OHLCData, OrderState } from "../../../types";
import { useState, useEffect, useCallback } from "react";
import { GetTradingSchedulesRequest, GetTradingSchedulesResponse, TradingSchedule } from "../../../../grpcGW/shares";
import { SeriesMarker, Time } from 'lightweight-charts';

type GetCandlesResponse = OHLCData[];

export const useGetCandles = (): (req: GetCandlesRequest) => Promise<GetCandlesResponse> => useIpcInvoke("GET_CANDLES");
export const useGetTradingSchedules = (): (req: GetTradingSchedulesRequest) => Promise<GetTradingSchedulesResponse> => useIpcInvoke("GET_TRADING_SCHEDULES");
export const useGetShares = () => useIpcInvoke("GET_SHARES");

// TODO: Нужен хук который сам бы хендлил отписку
const useSubscribeCandles = () => useIpcInvoke("SUBSCRIBE_CANDLES");
const useSubscribeOrders = () => useIpcInvoke("SUBSCRIBE_ORDER");
const useListenCandles = () => useIpcListen("NEW_CANDLE");
const useListenOrders = () => useIpcListen("NEW_ORDER");
// "BBG004730RP0" /* GAZP */
// "4c466956-d2ce-4a95-abb4-17947a65f18a" TGLD
// "BBG004730ZJ9" /* VTBR */
// "BBG004PYF2N3" /* POLY */

function orderToMarkerMapper(order: OrderState): SeriesMarker<Time> {
    return {
        time: order.time,
        position: order.operationType === 1 ? 'belowBar' : 'aboveBar',
        shape: 'circle',
        color: order.operationType === 1 ? 'green' : 'red',
        text: `${order.lotsExecuted} x ${order.price}`,
        size: 2,
    }
}
export const useOrders = (onNewOrder: (d: SeriesMarker<Time>) => void, figiOrInstrumentId = "4c466956-d2ce-4a95-abb4-17947a65f18a") => {
    const subscribe = useSubscribeOrders();
    const [registerOrderCb, unregisterOrderCb] = useListenOrders();

    const subscribeOrders = async () => {
        await subscribe({
            instrumentId: figiOrInstrumentId,
        });
    }

    const handleNewOrder = useCallback((e: Event, order: OrderState) => {
        if (!order) return;

        onNewOrder(orderToMarkerMapper(order));
    }, []);

    useEffect(() => {
        registerOrderCb(handleNewOrder);
    }, [handleNewOrder]);

    useEffect(() => {
        subscribeOrders();

        return () => {
            unregisterOrderCb(handleNewOrder);
        }
    }, [figiOrInstrumentId]);
}

export const useCandles = (onNewCandle: (d: OHLCData) => void, figiOrInstrumentId = "BBG004PYF2N3" /* POLY */, interval = 1) => {
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