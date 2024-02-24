import { useCallback, useEffect, useState } from "react";

// TODO: Заюзать useSyncExternalStorage для подписки на isAuthorised

export type AuthInfo =
    { isAuthorized: boolean; isLoaded: boolean; account: string; isSandbox?: boolean }

const DEFAULT_AUTH_INFO = { isAuthorized: false, isSandbox: true, account: '', isLoaded: false }

class AuthState {
    static instance: AuthState;
    constructor() {
        if (!!AuthState.instance) {
            return AuthState.instance;
        }

        AuthState.instance = this;

        return this;
    }

    public state: AuthInfo = DEFAULT_AUTH_INFO;
}

export const useAuth = () => {
    const authState = new AuthState();
    const [atuhInfo, setAuthInfo] = useState(authState.state);

    const getAuthInfo = useCallback(async () => {
        const info = await window.ipc.invoke('GET_AUTH_INFO');

        return info || { isAuthorised: false, isSandbox: true, account: null };
    }, []);

    const updateAuthInfo = async () => {
        try {
            const newAuthInfo = await getAuthInfo();
            authState.state = { isAuthorized: newAuthInfo.isAuthorised, isSandbox: newAuthInfo.isSandbox, account: newAuthInfo.account, isLoaded: true };
        } catch (e) {
            authState.state = { ...DEFAULT_AUTH_INFO, isLoaded: true };

            // TODO: Make notification
        }
        setAuthInfo(authState.state);
    }

    useEffect(() => {
        if (!authState.state.isLoaded) {
            updateAuthInfo()
        }
    }, []);

    return atuhInfo;
}