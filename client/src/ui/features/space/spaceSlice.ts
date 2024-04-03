import { createSlice } from '@reduxjs/toolkit';
import { Share } from '../../../../grpcGW/shares';

export type SpaceInfo = {
    shares: Share[]
};

export const initialState: SpaceInfo = {
    shares: []
};


const spaceSlice = createSlice({
    name: 'space',
    initialState,
    reducers: {
        setShares: (state, action) => {
            state.shares = action.payload;
        },
    },
});

export const {setShares} = spaceSlice.actions;

export default spaceSlice.reducer;