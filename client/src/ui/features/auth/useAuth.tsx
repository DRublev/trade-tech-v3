import { useCallback, useEffect, useState } from "react";

// TODO: Заюзать useSyncExternalStorage для подписки на isAuthorised

export type AuthInfo =
    { isAuthorized: false } |
    { isAuthorized: true, account: string, isSandbox?: boolean }

const DEFAULT_AUTH_INFO = { isAuthorized: false, isSandbox: true, account: '' }

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
console.log("31 useAuth", info);

        return info || { isAuthorised: false, isSandbox: true, account: null };
    }, []);

    const updateAuthInfo = async () => {
        try {
            const newAuthInfo = await getAuthInfo();
            authState.state = { isAuthorized: newAuthInfo.isAuthorised, isSandbox: newAuthInfo.isSandbox, account: newAuthInfo.account };
            setAuthInfo(authState.state);
        } catch (e) {
            authState.state = DEFAULT_AUTH_INFO;
            // TODO: Make notification
        }
    }

    useEffect(() => {
        updateAuthInfo()
    }, []);

    return atuhInfo;
}