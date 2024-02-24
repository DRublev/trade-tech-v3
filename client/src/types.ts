import { ipcEvents } from "./ipcEvents";

export type ValidChannel = keyof typeof ipcEvents;

export interface OHLCData {
    readonly close: number;
    readonly date: Date;
    readonly high: number;
    readonly low: number;
    readonly open: number;
    readonly volume: number;
}