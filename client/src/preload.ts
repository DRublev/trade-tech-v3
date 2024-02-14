// See the Electron documentation for details on how to use preload scripts:
// https://www.electronjs.org/docs/latest/tutorial/process-model#preload-scripts

import { ipcRenderer, contextBridge } from "electron";
import { ipcEvents } from "./ipcEvents";

const validChannels = Object.values(ipcEvents);

type ValidChannel = keyof typeof ipcEvents;

export interface IIpcRenderer {
    on: (channel: ValidChannel, listener: (event: any, ...args: any[]) => void) => void;
    once: (channel: ValidChannel, listener: (event: any, ...args: any[]) => void) => void;
    removeListener: (channel: ValidChannel, listener: (event: any, ...args: any[]) => void) => void;
    removeAllListeners: (channel: ValidChannel) => void;
    send: (channel: ValidChannel, ...args: any[]) => void;
    sendSync: (channel: ValidChannel, ...args: any[]) => void;
    sendToHost: (channel: ValidChannel, ...args: any[]) => void;
    invoke: (channel: ValidChannel, ...args: any[]) => Promise<any>;
}

export class SafeIpcRenderer {
    [x: string]: (channel: string, ...args: any[]) => any;
    constructor(events: string[]) {
        const protect = (fn: any) => {
            return (channel: string, ...args: any[]) => {
                if (!events.includes(channel)) {
                    throw new Error(`Blocked access to unknown channel ${channel} from the renderer. 
                          Add channel to whitelist in preload.js in case it is legitimate.`);
                }
                return fn.apply(ipcRenderer, [channel].concat(args));
            };
        };
        this.on = protect(ipcRenderer.on);
        this.once = protect(ipcRenderer.once);
        this.removeListener = protect(ipcRenderer.removeListener);
        this.removeAllListeners = protect(ipcRenderer.removeAllListeners);
        this.send = protect(ipcRenderer.send);
        this.sendSync = protect(ipcRenderer.sendSync);
        this.sendToHost = protect(ipcRenderer.sendToHost);
        this.invoke = protect(ipcRenderer.invoke);
    }
}

declare global {
    interface Window { ipc: IIpcRenderer; }
}

contextBridge.exposeInMainWorld(
    'ipc', new SafeIpcRenderer(validChannels),
);