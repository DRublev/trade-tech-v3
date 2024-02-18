import { useCallback } from "react";
import { ValidChannel } from "../../../types";


const useIpcInoke = (channel: ValidChannel) => {
    const invoke = useCallback((payload: unknown) => window.ipc ? window.ipc.invoke(channel, payload) : Promise.reject, []);

    return invoke;
};
export const useRegister = () => useIpcInoke("REGISTER");
export const useSetAccount = () => useIpcInoke("SET_ACCOUNT");
export const useGetAccount = () => useIpcInoke("GET_ACCOUNTS");
