import { useIpcInvoke } from '../../hooks';
import { useCallback } from 'react';
import { useAppDispatch } from '../../../store';
import { setRegisterData } from './registerSlice';



export const useRegister = <I, O>() => useIpcInvoke<I, O>("REGISTER");
export const useSetAccount = <I, O>() => useIpcInvoke<I, O>("SET_ACCOUNT");
export const useGetAccount = <I, O>() => useIpcInvoke<I, O>("GET_ACCOUNTS");

export const useRegistration = () => {
    const register = useRegister();
    const dispatch = useAppDispatch();

    const registerCallback = useCallback(async (data: Record<string, string>) => {
        dispatch(setRegisterData(data));
        await register(data);
    }, [register])

    return [registerCallback]
}