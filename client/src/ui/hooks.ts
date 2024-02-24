import { useCallback } from "react";
import { ValidChannel } from "../types";

export const useIpcInoke = <Req, Res>(channel: ValidChannel) => {
    const invoke: (r: Req) => Promise<Res> = useCallback((payload: Req) => (window.ipc ? window.ipc.invoke(channel, payload) : Promise.reject) as any, []);

    return invoke;
};

export const useIpcListen = (channel: ValidChannel) => {
    const listen = useCallback((cb: any) => window.ipc ? window.ipc.on(channel, cb) : Promise.reject, [channel]);
    const removeListen = useCallback((cb: any) => window.ipc ? window.ipc.removeListener(channel, cb) : Promise.reject, [channel]);

    return [listen, removeListen];
}