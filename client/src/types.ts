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