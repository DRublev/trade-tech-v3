import { createSlice } from '@reduxjs/toolkit';
import { Share } from '../../../node/grpc/contracts/shares';


export type SpaceInfo = {
    shares: Share[];
    currentInstrument: string | null;
    initiallySetCurrentInstrument: boolean;
};

export const initialState: SpaceInfo = {
    shares: [],
    currentInstrument: null,
    initiallySetCurrentInstrument: false,
};


const spaceSlice = createSlice({
    name: 'space',
    initialState,
    reducers: {
        setShares: (state, action) => {
            state.shares = action.payload;
        },
        setCurrentInstrument: (state, action) => {
            state.currentInstrument = action.payload;
        },
        setInitiallySetCurrentInstrument: (state, action) => {
            state.initiallySetCurrentInstrument = action.payload;
        },
    },
});

export const { setShares, setCurrentInstrument, setInitiallySetCurrentInstrument } = spaceSlice.actions;

export default spaceSlice.reducer;