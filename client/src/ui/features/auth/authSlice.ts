import { createSlice } from '@reduxjs/toolkit';

export type AuthInfo =
    {
        isAuthorized: boolean;
        isLoaded: boolean;
        account: string;
        isSandbox?: boolean;
        error: string | null;
    }

export const DEFAULT_AUTH_INFO: AuthInfo = {
    error: '',
    isAuthorized: false,
    isSandbox: true,
    account: '',
    isLoaded: false
};


const authSlice = createSlice({
    name: 'auth',
    initialState: DEFAULT_AUTH_INFO,
    reducers: {
        setCurrentAccount: (state, {payload}) => {
          state.account = payload.account;
        },
        setAuthData: (state, {payload}) => {
            state.isLoaded = payload.isLoaded;
            state.isSandbox = payload.isSandbox;
            state.isAuthorized = payload.isAuthorized;
        },
        setLoaded: (state, {payload}) => {
          state.isLoaded = payload.isLoaded;
        },
        setError: (state, {payload}) => {
            state.error = payload.error;
            state.account = null;
            state.isLoaded = true;
            state.isSandbox = true;
            state.isAuthorized = false;
        },
        logout: (state) => {
            state.isAuthorized = false;
            state.account = null;
        },
    },
});

export const {logout, setAuthData, setError, setLoaded, setCurrentAccount} = authSlice.actions;

export default authSlice.reducer;