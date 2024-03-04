import { GetCandlesRequest } from "../../../../grpcGW/marketData";
import { useIpcInoke, useIpcListen } from "../../hooks";
import { OHLCData } from "../../../types";
import { useState, useEffect, useCallback } from "react";
import { GetInstrumentsRequest } from "../../../../grpcGW/shares";
import { GetSharesResponse } from "../../../../grpcGW/shares";

type GetCandlesResponse = OHLCData[];

export const useGetCandles = (): (req: GetCandlesRequest) => Promise<GetCandlesResponse> => useIpcInoke("GET_CANDLES");
// export const useGetShares = (): (req: GetInstrumentsRequest) => Promise<GetSharesResponse> => useIpcInoke("GET_SHARES");
export const useGetShares = () => useIpcInoke("GET_SHARES");

// TODO: Нужен хук который сам бы хендлил отписку
const useSubscribeCandles = () => useIpcInoke("SUBSCRIBE_CANDLES");
const useListenCandles = () => useIpcListen("NEW_CANDLE");

export const useCandles = (figiOrInstrumentId = "4c466956-d2ce-4a95-abb4-17947a65f18a" /* TGLD */, interval = 1) => {
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

            console.log("39 hooks", candles);

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

    const handleNewCandle = useCallback((e: Event, candle: OHLCData) => {
        if (!candle) return;

        setData((prevData) => {
            if (!prevData.length) return [candle];

            const lastCandleDate = prevData[prevData.length - 1].date;

            if (lastCandleDate.getMinutes() === candle.date.getMinutes()) {
                prevData[prevData.length - 1] = candle;
                return [...prevData];
            }

            return [...prevData, candle];
        });
        
    }, [figiOrInstrumentId])

    useEffect(() => {
        onCandles(handleNewCandle);



    }, [handleNewCandle]);

    useEffect(() => {
        console.log("86 hooks", );
        
        getInitialCandels();
        subscribeCandles();

        return () => {
            off(handleNewCandle);
        }
    }, [figiOrInstrumentId]);

    return { data, isLoading, error };
}