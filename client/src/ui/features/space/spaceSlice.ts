import { createSlice } from '@reduxjs/toolkit';
import { Share } from '../../../node/grpc/contracts/shares';


// "BBG004730RP0" /* GAZP */
// "4c466956-d2ce-4a95-abb4-17947a65f18a" TGLD
// "BBG004730ZJ9" /* VTBR */
// "BBG004PYF2N3" /* POLY */
const defaultInstrumentId = "BBG004PYF2N3" /* POLY */;

export type SpaceInfo = {
    shares: Share[];
    currentInstrument: string;
};

export const initialState: SpaceInfo = {
    shares: [],
    currentInstrument: defaultInstrumentId,
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
        }
    },
});

export const { setShares, setCurrentInstrument } = spaceSlice.actions;

export default spaceSlice.reducer;