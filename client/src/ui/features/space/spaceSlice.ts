import { createSlice } from '@reduxjs/toolkit';
import { Share } from '../../../node/grpc/contracts/shares';
import type { StrategyKey } from '../strategy/types';


export type SpaceInfo = {
    shares: Share[];
    currentInstrument: string | null;
    strategy: StrategyKey;
    initiallySetCurrentInstrument: boolean;
};

export const initialState: SpaceInfo = {
    shares: [],
    currentInstrument: null,
    strategy: 'spread_v0',
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
        setStrategy: (state, action) => {
            state.strategy = action.payload;
        },
        setInitiallySetCurrentInstrument: (state, action) => {
            state.initiallySetCurrentInstrument = action.payload;
        },
    },
});

export const { setShares, setCurrentInstrument, setInitiallySetCurrentInstrument, setStrategy } = spaceSlice.actions;

export default spaceSlice.reducer;