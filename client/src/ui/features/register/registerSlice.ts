import { createSlice } from '@reduxjs/toolkit';

export type RegisterInfo = {
    registerInfo: Record<string, string>
};

//TODO Нормально типизировать
export const initialState: RegisterInfo = {
    registerInfo: {}
};


const registerSlice = createSlice({
    name: 'register',
    initialState,
    reducers: {
        setRegisterData: (state, action) => {
            state.registerInfo = action.payload;
        },
    },
});

export const {setRegisterData} = registerSlice.actions;

export default registerSlice.reducer;