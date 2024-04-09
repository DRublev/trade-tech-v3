import { useCallback, useEffect } from "react";
import { useAppDispatch, useAppSelector } from '../../../store';
import { DEFAULT_AUTH_INFO, setAuthData, setLoaded } from './authSlice';

// TODO: Заюзать useSyncExternalStorage для подписки на isAuthorised

export type AuthInfo =
    { isAuthorized: boolean; isLoaded: boolean; account: string; isSandbox?: boolean }

export const useAuth = () => {
    const authState = useAppSelector(state => state.auth)
    const isLoaded = useAppSelector(state => state.auth.isLoaded)
    const dispatch = useAppDispatch();

    const getAuthInfo = useCallback(async () => {
        const info = await window.ipc.invoke('GET_AUTH_INFO');

        return info || { isAuthorised: false, isSandbox: false, account: null };
    }, []);

    const setShouldUpdateAuthInfo = useCallback(() => {
        dispatch(setLoaded(false))
    }, []);

    const updateAuthInfo = async () => {
        try {
            const newAuthInfo = await getAuthInfo();
            dispatch(setAuthData({ isAuthorized: newAuthInfo.isAuthorised, isSandbox: newAuthInfo.isSandbox, account: newAuthInfo.account, isLoaded: true }))
        } catch (e) {
            dispatch(setAuthData({ ...DEFAULT_AUTH_INFO, isLoaded: true }))

            // TODO: Make notification
        }
    }

    useEffect(() => {
        if (!isLoaded) {
            updateAuthInfo()
        }
    }, [isLoaded]);

    return { ...authState, setShouldUpdateAuthInfo };
}