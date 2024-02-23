import { useCallback } from "react";
import { ValidChannel } from "../types";

export const useIpcInoke = <Req, Res>(channel: ValidChannel) => {
    const invoke: (r: Req) => Promise<Res> = useCallback((payload: Req) => (window.ipc ? window.ipc.invoke(channel, payload) : Promise.reject) as any, []);

    return invoke;
};