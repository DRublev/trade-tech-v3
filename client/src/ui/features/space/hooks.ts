import { GetCandlesRequest } from "../../../../grpcGW/marketData";
import { useIpcInoke } from "../../hooks";
import { OHLCData } from "../../../types";
import { useState, useEffect } from "react";

type GetCandlesResponse = OHLCData[];

export const useGetCandles = (): (req: GetCandlesRequest) => Promise<GetCandlesResponse> => useIpcInoke("GET_CANDLES");

export const useHistoricCandels = (figiOrInstrumentId: string = "BBG004730N88") => {
    const getCandles = useGetCandles();
    const [data, setData] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);

    const now = new Date();
    const dayAgoDate = new Date(new Date(now).setDate(now.getDate() - 1))

    const getInitialCandels = async () => {
        try {
            setIsLoading(true);

            const candles = await getCandles({
                instrumentId: figiOrInstrumentId,
                interval: 1,
                start: dayAgoDate,
                end: now,
            });
            setData(candles.filter(d => d));
        } catch (e) {
            setError(e);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        getInitialCandels();
    }, [figiOrInstrumentId]);

    return { data, isLoading, error };
}