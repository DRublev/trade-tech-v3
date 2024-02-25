import { GetCandlesRequest } from "../../../../grpcGW/marketData";
import { useIpcInoke, useIpcListen } from "../../hooks";
import { OHLCData } from "../../../types";
import { useState, useEffect, useCallback } from "react";
import { GetInstrumentsRequest } from "../../../../grpcGW/instruments";
import { GetSharesResponse } from "../../../../grpcGW/instruments";

type GetCandlesResponse = OHLCData[];

export const useGetCandles = (): (req: GetCandlesRequest) => Promise<GetCandlesResponse> => useIpcInoke("GET_CANDLES");
// export const useGetShares = (): (req: GetInstrumentsRequest) => Promise<GetSharesResponse> => useIpcInoke("GET_SHARES");
export const useGetShares = () => useIpcInoke("GET_SHARES");

// TODO: Нужен хук который сам бы хендлил отписку
const useSubscribeCandles = () => useIpcInoke("SUBSCRIBE_CANDLES");
const useListenCandles = () => useIpcListen("NEW_CANDLE");

export const useCandles = (figiOrInstrumentId: string = "BBG004730N88", interval = 1) => {
    const getCandles = useGetCandles();
    const subscribe = useSubscribeCandles();

    const [onCandles, off] = useListenCandles();

    const [data, setData] = useState<OHLCData[]>([]);
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
            setData(candles.filter(d => d));
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
    const sub = async () => {
        await onCandles(handleNewCandle)
    };

    const handleNewCandle = (e: Event, candle: OHLCData) => {
        console.log("56 hooks", candle);
        if (!candle) return;
        const lastCandleDate = data[data.length - 1].date;
        if (lastCandleDate.getMinutes() === candle.date.getMinutes()) {
            setData(data.splice(data.length - 1, 1, candle));
        } else {
            data.push(candle);
            setData(data);
        }
    }

    useEffect(() => {
        getInitialCandels();
        sub();
        subscribeCandles();

        return () => {
            off(handleNewCandle);
        }

    }, [figiOrInstrumentId]);

    return { data, isLoading, error };
}