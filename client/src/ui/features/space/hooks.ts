import { GetCandlesRequest } from "../../../../grpcGW/marketData";
import { useIpcInoke } from "../../hooks";
import { OHLCData } from "../../../types";

type GetCandlesResponse = OHLCData[];

export const useGetCandles = (): (req: GetCandlesRequest) => Promise<GetCandlesResponse> => useIpcInoke("GET_CANDLES");