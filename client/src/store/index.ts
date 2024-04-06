import { configureStore } from '@reduxjs/toolkit';
import accountsSlice from '../ui/features/accounts/accountsSlice';
import authSlice from '../ui/features/auth/authSlice';
import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';
import registerSlice from '../ui/features/register/registerSlice';
import spaceSlice from '../ui/features/space/spaceSlice';

export const store = configureStore({
    reducer: {
        accounts: accountsSlice,
        auth: authSlice,
        register: registerSlice,
        space: spaceSlice,
    },
});
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;