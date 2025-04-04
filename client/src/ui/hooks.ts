import { Bindings, ChildLoggerOptions } from "pino";
import { useCallback } from "react";
import logger from "../logger";
import { ValidChannel } from "../types";

export const useIpcInvoke = <Req, Res>(channel: ValidChannel) => {
    const invoke: (r: Req) => Promise<Res> = useCallback((payload: Req) => (window.ipc ? window.ipc.invoke(channel, payload) : Promise.reject) as any, []);

    return invoke;
};

// TODO: FIXME: Есть бага с тем, что мемоизирутеся коллбек и у него в области видимости не актуальные переменные (не подтягивается актуальный стейт компонента)
export const useIpcListen = (channel: ValidChannel) => {
    const listen = useCallback((cb: any) => window.ipc ? window.ipc.on(channel, cb) : Promise.reject, [channel]);
    const removeListen = useCallback((cb: any) => window.ipc ? window.ipc.removeListener(channel, cb) : Promise.reject, [channel]);

    return [listen, removeListen];
};

type TypedBindings = Bindings & {
    component: string;
};

export const useLogger = (bindings?: TypedBindings, options?: ChildLoggerOptions) => {
    if (bindings) {
        return logger.child(bindings, options);
    }

    return logger;
};
