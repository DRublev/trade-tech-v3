import { createSlice } from '@reduxjs/toolkit';
import { Share } from '../../../node/grpc/contracts/shares';


// "BBG004730RP0" /* GAZP */
// "4c466956-d2ce-4a95-abb4-17947a65f18a" /* TGLD */
// "BBG004730ZJ9" /* VTBR */
// "BBG004PYF2N3" /* POLY */
// "b71bd174-c72c-41b0-a66f-5f9073e0d1f5" /* VKCO */
// "4163e41d-55f4-4f93-82fc-6c44fe5d444e" /* SPBE */
// "BBG00F9XX7H4" /* RNFT */
// "9654c2dd-6993-427e-80fa-04e80a1cf4da" /* TMOS */
const defaultInstrumentId = "4c466956-d2ce-4a95-abb4-17947a65f18a" /* TGLD */;

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