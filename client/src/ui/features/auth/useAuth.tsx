import { useCallback } from "react";
import { useAppDispatch } from '../../../store';
import { useGetAccount, useSetAccount } from "../register/hooks";
import { DEFAULT_AUTH_INFO, setAuthData, setCurrentAccount } from './authSlice';

// TODO: Заюзать useSyncExternalStorage для подписки на isAuthorised

export type AuthInfo =
    { isAuthorized: boolean; isLoaded: boolean; account: string; isSandbox?: boolean }

export const useAuth = () => {
    const setAccount = useSetAccount();
    const dispatch = useAppDispatch();
    const getAccount = useGetAccount();
  
    const getAuthInfo = useCallback(async () => {
        const info = await window.ipc.invoke('GET_AUTH_INFO');

        return info || { isAuthorised: false, isSandbox: false, account: null };
    }, []);

    const updateAuthInfo = async () => {
        try {
            const newAuthInfo = await getAuthInfo();
            
            dispatch(
              setAuthData({ isAuthorized: newAuthInfo.isAuthorised, isSandbox: newAuthInfo.isSandbox, isLoaded: true })
            )
        } catch (e) {
            dispatch(setAuthData({ ...DEFAULT_AUTH_INFO, isLoaded: true }))

            // TODO: Make notification
        }
    }

    const updateAuth = useCallback(async () => {
      await updateAuthInfo();
      const { AccountId } = await getAccount({});
      dispatch(setCurrentAccount({ account: AccountId }))
    }, [setAccount ,updateAuthInfo]);

    const selectAccount = useCallback(async (id: FormDataEntryValue) => {
      await setAccount({ id });
      dispatch(setCurrentAccount({account: id}))
      await updateAuth();
    }, [setAccount ,updateAuthInfo]);

    return { selectAccount, updateAuthInfo, updateAuth };
}