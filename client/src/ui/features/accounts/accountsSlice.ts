import { createSlice } from '@reduxjs/toolkit';

export type Account = {
    id: string;
    name: string;
}

export type RawAccount = {
    Id: string;
    Name: string;
}

export type AccountsInfo =
    {
        accounts: Account[];
        error: string | null;
        selectedAccountId: string;
    }

export const initialState: AccountsInfo = {
    accounts: [],
    selectedAccountId: null,
    error: null
};


const accountsSlice = createSlice({
    name: 'accounts',
    initialState,
    reducers: {
        setAccounts: (state, action) => {
            state.accounts = action.payload;
        },
    },
});

export const {setAccounts} = accountsSlice.actions;

export default accountsSlice.reducer;