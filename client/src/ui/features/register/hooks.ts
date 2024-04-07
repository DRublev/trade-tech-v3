import { ipcEvents } from "../../../ipcEvents";
import { useIpcInvoke } from "../../hooks";
import { useCallback } from 'react';
import { useAppDispatch } from '../../../store';
import { setRegisterData } from './registerSlice';

export const useRegister = () => useIpcInvoke(ipcEvents.REGISTER);
export const usePruneTokens = () => useIpcInvoke(ipcEvents.PRUNE_TOKENS);
export const useSetAccount = () => useIpcInvoke(ipcEvents.SET_ACCOUNT);
export const useGetAccount = () => useIpcInvoke(ipcEvents.GET_ACCOUNTS);

export const useRegistration = () => {
    const register = useRegister();
    const dispatch = useAppDispatch();

    const registerCallback = useCallback(async (data: Record<string, string>) => {
        dispatch(setRegisterData(data));
        await register(data);
    }, [register])

    return [registerCallback]
}
