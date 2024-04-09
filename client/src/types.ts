import { UTCTimestamp } from "lightweight-charts";
import { ipcEvents } from "./ipcEvents";

export type ValidChannel = keyof typeof ipcEvents;

export interface OHLCData {
    /**
     * @example 1529899200 - Literal timestamp representing 2018-06-25T04:00:00.000Z
     */
    readonly time: UTCTimestamp;
    readonly open: number;
    readonly high: number;
    readonly low: number;
    readonly close: number;
    readonly volume: number;
}

export enum OrderOperations {
    Buy = 1,
    Sell = 2
}

export interface OrderState {
    readonly id: string;
    readonly instrumentId: string;
    readonly price: number
    readonly status: number;
    readonly lotsRequested: number;
    readonly lotsExecuted: number;
    readonly operationType: OrderOperations;
    readonly time: UTCTimestamp;
    readonly strategy: string;
}