import { ipcEvents } from "./ipcEvents";

export type ValidChannel = keyof typeof ipcEvents;
